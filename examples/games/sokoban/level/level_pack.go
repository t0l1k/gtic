package level

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gtic/examples/games/sokoban/assets"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type LevelPack struct {
	Name   string
	File   string
	Levels []string
}

func LoadPacks() ([]LevelPack, error) {
	return LoadPacksFrom(assets.LevelFiles, "levels")
}

func LoadPacksFrom(source fs.FS, dir string) ([]LevelPack, error) {
	entries, err := fs.ReadDir(source, dir)
	if err != nil {
		return nil, fmt.Errorf("read level collection directory %q: %w", dir, err)
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".xsb" || ext == ".txt" {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)
	if len(files) == 0 {
		return nil, fmt.Errorf("no .xsb or .txt level collections found in assets/levels")
	}

	packs := make([]LevelPack, 0, len(files))
	for _, file := range files {
		levelPath := path.Join(dir, file)
		data, err := fs.ReadFile(source, levelPath)
		if err != nil {
			return nil, fmt.Errorf("read level pack %q: %w", file, err)
		}
		levels, err := ParseLevels(string(data))
		if err != nil {
			return nil, fmt.Errorf("parse level pack %q: %w", file, err)
		}
		packs = append(packs, LevelPack{
			Name:   strings.TrimSuffix(file, filepath.Ext(file)),
			File:   file,
			Levels: levels,
		})
	}
	return packs, nil
}

func ParseLevels(data string) ([]string, error) {
	var (
		levels []string
		rows   []string
	)
	metadata := false
	flush := func() error {
		if len(rows) == 0 {
			return nil
		}
		width := 0
		playerCount := 0
		for _, row := range rows {
			width = max(width, len(row))
			for _, tile := range row {
				if tile == '@' || tile == '+' {
					playerCount++
				}
			}
		}
		if playerCount != 1 {
			return fmt.Errorf("level %d has %d players; want exactly one", len(levels)+1, playerCount)
		}
		var level strings.Builder
		for i, row := range rows {
			if i > 0 {
				level.WriteByte('\n')
			}
			level.WriteString(normalizeLevelRow(row))
			level.WriteString(strings.Repeat(" ", width-len(row)))
		}
		levels = append(levels, level.String())
		rows = nil
		return nil
	}

	scanner := bufio.NewScanner(strings.NewReader(data))
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		line = strings.TrimSuffix(line, "\r")
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, ";") {
			if err := flush(); err != nil {
				return nil, err
			}
			metadata = false
			continue
		}
		if IsLevelNumber(trimmed) {
			if err := flush(); err != nil {
				return nil, err
			}
			metadata = true
			continue
		}
		if trimmed == "" {
			if err := flush(); err != nil {
				return nil, err
			}
			continue
		}
		if IsCollectionMetadata(trimmed) {
			if err := flush(); err != nil {
				return nil, err
			}
			metadata = true
			continue
		}
		if metadata && !IsLevelRow(line) {
			continue
		}
		if !IsLevelRow(line) {
			if err := flush(); err != nil {
				return nil, err
			}
			metadata = true
			continue
		}
		metadata = false
		rows = append(rows, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if err := flush(); err != nil {
		return nil, err
	}
	if len(levels) == 0 {
		return nil, fmt.Errorf("no levels found")
	}
	return levels, nil
}

func IsCollectionMetadata(line string) bool {
	lower := strings.ToLower(line)
	for _, field := range []string{"title:", "description:", "author:", "email:", "website:", "homepage:", "set:", "copyright:"} {
		if strings.HasPrefix(lower, field) {
			return true
		}
	}
	return false
}

func IsLevelNumber(line string) bool {
	if line == "" {
		return false
	}
	for _, char := range line {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func IsLevelRow(line string) bool {
	for _, tile := range line {
		switch tile {
		case ' ', '#', '.', '@', '+', '$', '*', '-', '_':
		default:
			return false
		}
	}
	return len(line) > 0
}

func normalizeLevelRow(row string) string {
	return strings.NewReplacer("-", " ", "_", " ").Replace(row)
}

type SavedProgress map[string]int

type AttemptResult struct {
	Attempt int `json:"attempt"`
	Moves   int `json:"moves"`
	Pushes  int `json:"pushes"`
}

type LevelStatistics struct {
	Attempts  int             `json:"attempts"`
	Successes int             `json:"successes"`
	Best      *AttemptResult  `json:"best,omitempty"`
	Results   []AttemptResult `json:"results,omitempty"`
}

type SavedStatistics map[string]map[int]LevelStatistics

func LoadStatistics(path string) (SavedStatistics, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return SavedStatistics{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read Sokoban statistics: %w", err)
	}
	var stats SavedStatistics
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("decode Sokoban statistics: %w", err)
	}
	if stats == nil {
		stats = SavedStatistics{}
	}
	return stats, nil
}

func SaveStatistics(path string, stats SavedStatistics) error {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("encode Sokoban statistics: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create statistics directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write Sokoban statistics: %w", err)
	}
	return nil
}

func (s SavedStatistics) Level(pack string, level int) LevelStatistics {
	return s[pack][level]
}

func (s SavedStatistics) StartAttempt(pack string, level int) int {
	if s[pack] == nil {
		s[pack] = make(map[int]LevelStatistics)
	}
	stats := s[pack][level]
	stats.Attempts++
	s[pack][level] = stats
	return stats.Attempts
}

func (s SavedStatistics) RecordSuccess(pack string, level int, result AttemptResult) {
	if s[pack] == nil {
		s[pack] = make(map[int]LevelStatistics)
	}
	stats := s[pack][level]
	stats.Successes++
	stats.Results = append(stats.Results, result)
	if stats.Best == nil || result.Moves < stats.Best.Moves ||
		(result.Moves == stats.Best.Moves && result.Pushes < stats.Best.Pushes) {
		best := result
		stats.Best = &best
	}
	s[pack][level] = stats
}

func LoadProgress() (SavedProgress, string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, "", fmt.Errorf("find user config directory: %w", err)
	}
	path := filepath.Join(configDir, "gtic", "sokoban-progress.json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return SavedProgress{}, path, nil
	}
	if err != nil {
		return nil, "", fmt.Errorf("read Sokoban progress: %w", err)
	}
	var progress SavedProgress
	if err := json.Unmarshal(data, &progress); err != nil {
		return nil, "", fmt.Errorf("decode Sokoban progress: %w", err)
	}
	if progress == nil {
		progress = SavedProgress{}
	}
	return progress, path, nil
}

func SaveProgress(path string, progress SavedProgress) error {
	data, err := json.MarshalIndent(progress, "", "  ")
	if err != nil {
		return fmt.Errorf("encode Sokoban progress: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create progress directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write Sokoban progress: %w", err)
	}
	return nil
}
