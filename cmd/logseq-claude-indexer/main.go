package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/dyluth/logseq-claude-indexer/internal/indexer"
	"github.com/dyluth/logseq-claude-indexer/internal/parser"
	"github.com/dyluth/logseq-claude-indexer/internal/scanner"
	"github.com/dyluth/logseq-claude-indexer/internal/writer"
	"github.com/dyluth/logseq-claude-indexer/pkg/models"
)

var (
	repoPath  string
	outputDir string
	quiet     bool
	verbose   bool
	dryRun    bool
	version   = "0.1.0"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "logseq-claude-indexer",
	Short: "Generate Claude Code-optimized indexes from Logseq repositories",
	Long: `logseq-claude-indexer scans Logseq markdown files and generates
structured index files for AI consumption. It extracts tasks, time tracking,
and page references to help Claude Code understand your knowledge base.`,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate all indexes",
	Long:  `Scan the Logseq repository and generate both task index and reference graph.`,
	RunE:  runGenerate,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("logseq-claude-indexer version %s\n", version)
	},
}

func init() {
	// Add commands
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(versionCmd)

	// Add flags to generate command
	generateCmd.Flags().StringVar(&repoPath, "repo", ".", "Path to Logseq repository")
	generateCmd.Flags().StringVar(&outputDir, "output", ".claude/indexes", "Output directory for index files")
	generateCmd.Flags().BoolVar(&quiet, "quiet", false, "Suppress output (for git hooks)")
	generateCmd.Flags().BoolVar(&verbose, "verbose", false, "Show detailed logging")
	generateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be generated without writing files")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Set up logger
	logger := log.New(os.Stdout, "", 0)
	if quiet {
		logger.SetOutput(io.Discard)
	}

	logger.Printf("Scanning Logseq repository: %s", repoPath)

	// Convert to absolute path
	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return fmt.Errorf("invalid repo path: %w", err)
	}

	// Verify repo path exists
	if _, err := os.Stat(absRepoPath); os.IsNotExist(err) {
		return fmt.Errorf("repository path does not exist: %s", absRepoPath)
	}

	// 1. Scan for files
	if verbose {
		logger.Println("Step 1: Scanning for markdown files...")
	}
	sc := scanner.New(absRepoPath)
	files, err := sc.Scan()
	if err != nil {
		return fmt.Errorf("scanning files: %w", err)
	}

	logger.Printf("Found %d markdown files", len(files))

	if len(files) == 0 {
		logger.Println("No markdown files found in pages/ or journals/")
		return nil
	}

	// 2. Parse all files
	if verbose {
		logger.Println("Step 2: Parsing files...")
	}

	var allTasks []models.Task
	var allRefs []models.PageReference
	parseErrors := 0

	for _, file := range files {
		content, err := os.ReadFile(file.AbsolutePath)
		if err != nil {
			if verbose {
				logger.Printf("Warning: Failed to read %s: %v", file.Path, err)
			}
			parseErrors++
			continue
		}

		// Parse tasks
		tasks, err := parser.ParseTasks(string(content), file.Path)
		if err != nil {
			if verbose {
				logger.Printf("Warning: Failed to parse tasks in %s: %v", file.Path, err)
			}
			parseErrors++
		} else {
			allTasks = append(allTasks, tasks...)
		}

		// Parse references
		refs, err := parser.ParseReferences(string(content), file.Path)
		if err != nil {
			if verbose {
				logger.Printf("Warning: Failed to parse references in %s: %v", file.Path, err)
			}
			parseErrors++
		} else {
			allRefs = append(allRefs, refs...)
		}
	}

	logger.Printf("Extracted %d tasks and %d references", len(allTasks), len(allRefs))

	if parseErrors > 0 {
		logger.Printf("Warning: %d parse errors encountered", parseErrors)
	}

	// 3. Build indexes
	if verbose {
		logger.Println("Step 3: Building indexes...")
	}

	taskIndex := indexer.BuildTaskIndex(allTasks)
	graphIndex := indexer.BuildReferenceGraph(allRefs, files)
	timelineIndex := indexer.BuildTimelineIndex(allTasks, files)
	missingPagesIndex := indexer.BuildMissingPagesIndex(graphIndex, 5)

	if dryRun {
		logger.Println("\n=== DRY RUN MODE ===")
		logger.Printf("Would create task index with %d tasks", taskIndex.TotalTasks)
		logger.Printf("Would create reference graph with %d nodes", len(graphIndex.Nodes))
		logger.Printf("Would create timeline with %d days", len(timelineIndex.Entries))
		logger.Printf("Would create missing pages report with %d pages", len(missingPagesIndex.MissingPages))
		return nil
	}

	// 4. Write output files
	if verbose {
		logger.Println("Step 4: Writing index files...")
	}

	// Make output path absolute (relative to repo path)
	absOutputDir := outputDir
	if !filepath.IsAbs(outputDir) {
		absOutputDir = filepath.Join(absRepoPath, outputDir)
	}

	// Write task index by status
	if err := writer.WriteTaskIndex(taskIndex, absOutputDir); err != nil {
		return fmt.Errorf("writing task index: %w", err)
	}
	logger.Printf("✓ Created %s", filepath.Join(absOutputDir, "tasks-by-status.md"))

	// Write priority index
	if err := writer.WritePriorityIndex(taskIndex, absOutputDir); err != nil {
		return fmt.Errorf("writing priority index: %w", err)
	}
	logger.Printf("✓ Created %s", filepath.Join(absOutputDir, "tasks-by-priority.md"))

	// Write timeline recent
	if err := writer.WriteTimelineRecent(timelineIndex, absOutputDir); err != nil {
		return fmt.Errorf("writing recent timeline: %w", err)
	}
	logger.Printf("✓ Created %s", filepath.Join(absOutputDir, "timeline-recent.md"))

	// Write timeline full
	if err := writer.WriteTimelineFull(timelineIndex, absOutputDir); err != nil {
		return fmt.Errorf("writing full timeline: %w", err)
	}
	logger.Printf("✓ Created %s", filepath.Join(absOutputDir, "timeline-full.md"))

	// Write missing pages
	if err := writer.WriteMissingPages(missingPagesIndex, absOutputDir); err != nil {
		return fmt.Errorf("writing missing pages: %w", err)
	}
	logger.Printf("✓ Created %s", filepath.Join(absOutputDir, "missing-pages.md"))

	// Write reference graph
	if err := writer.WriteReferenceGraph(graphIndex, absOutputDir); err != nil {
		return fmt.Errorf("writing reference graph: %w", err)
	}
	logger.Printf("✓ Created %s", filepath.Join(absOutputDir, "reference-graph.md"))

	logger.Println("Index generation complete!")

	return nil
}
