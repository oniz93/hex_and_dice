package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all server configuration.
type Config struct {
	Port                 string        `json:"port"`
	RedisURL             string        `json:"redis_url"`
	LogLevel             string        `json:"log_level"`
	CORSOrigins          []string      `json:"cors_origins"`
	BalanceFile          string        `json:"balance_file"`
	WSPingInterval       time.Duration `json:"ws_ping_interval"`
	WSPongTimeout        time.Duration `json:"ws_pong_timeout"`
	ReconnectTimeout     time.Duration `json:"reconnect_timeout"`
	RoomTTL              time.Duration `json:"room_ttl"`
	ShutdownDrainTimeout time.Duration `json:"shutdown_drain_timeout"`
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Port:                 envOrDefault("PORT", "8080"),
		RedisURL:             envOrDefault("REDIS_URL", "redis://localhost:6379"),
		LogLevel:             envOrDefault("LOG_LEVEL", "INFO"),
		CORSOrigins:          strings.Split(envOrDefault("CORS_ORIGINS", "*"), ","),
		BalanceFile:          envOrDefault("BALANCE_FILE", "data/balance.yaml"),
		WSPingInterval:       durationOrDefault("WS_PING_INTERVAL", 15*time.Second),
		WSPongTimeout:        durationOrDefault("WS_PONG_TIMEOUT", 10*time.Second),
		ReconnectTimeout:     durationOrDefault("RECONNECT_TIMEOUT", 60*time.Second),
		RoomTTL:              durationOrDefault("ROOM_TTL", 5*time.Minute),
		ShutdownDrainTimeout: durationOrDefault("SHUTDOWN_DRAIN_TIMEOUT", 30*time.Second),
	}
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func durationOrDefault(key string, defaultVal time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return defaultVal
}

// BalanceData holds all game balance constants loaded from YAML.
type BalanceData struct {
	Economy     EconomyConfig              `yaml:"economy"`
	Troops      map[string]TroopConfig     `yaml:"troops"`
	Structures  map[string]StructureConfig `yaml:"structures"`
	NeutralMod  NeutralModConfig           `yaml:"neutral_modifiers"`
	Terrain     map[string]TerrainConfig   `yaml:"terrain"`
	Healing     HealingConfig              `yaml:"healing"`
	SuddenDeath SuddenDeathConfig          `yaml:"sudden_death"`
	MapGen      MapGenConfig               `yaml:"map_generation"`
	Matchmaking MatchmakingConfig          `yaml:"matchmaking"`
	WinCond     WinCondConfig              `yaml:"win_conditions"`
}

type EconomyConfig struct {
	StartingCoins   int `yaml:"starting_coins"`
	PassiveIncome   int `yaml:"passive_income"`
	StructureIncome int `yaml:"structure_income"`
}

type TroopConfig struct {
	Cost                    int    `yaml:"cost"`
	HP                      int    `yaml:"hp"`
	ATK                     int    `yaml:"atk"`
	DEF                     int    `yaml:"def"`
	Mobility                int    `yaml:"mobility"`
	Range                   int    `yaml:"range"`
	Damage                  string `yaml:"damage"`
	AntiStructureMultiplier int    `yaml:"anti_structure_multiplier,omitempty"`
}

type StructureConfig struct {
	HP     int    `yaml:"hp"`
	ATK    int    `yaml:"atk"`
	DEF    int    `yaml:"def"`
	Range  int    `yaml:"range"`
	Damage string `yaml:"damage"`
	Income int    `yaml:"income"`
	Spawn  bool   `yaml:"spawn"`
}

type NeutralModConfig struct {
	ATKReduction   int `yaml:"atk_reduction"`
	DamageStepDown int `yaml:"damage_step_down"`
}

type TerrainConfig struct {
	MovementCost int   `yaml:"movement_cost,omitempty"`
	ATKModifier  int   `yaml:"atk_modifier,omitempty"`
	DEFModifier  int   `yaml:"def_modifier,omitempty"`
	Passable     *bool `yaml:"passable,omitempty"` // nil means passable (default true)
}

type HealingConfig struct {
	PassiveRate int `yaml:"passive_rate"`
}

type SuddenDeathConfig struct {
	TurnThresholds map[string]int `yaml:"turn_thresholds"`
	ShrinkRate     int            `yaml:"shrink_rate"`
}

type MapGenConfig struct {
	NoiseThresholds  map[string]float64 `yaml:"noise_thresholds"`
	StructureCounts  map[string]int     `yaml:"structure_counts"`
	MinPassableRatio float64            `yaml:"min_passable_ratio"`
	MaxRetries       int                `yaml:"max_retries"`
}

type MatchmakingConfig struct {
	QuickMatchDefaults QuickMatchDefaults `yaml:"quick_match_defaults"`
}

type QuickMatchDefaults struct {
	MapSize   string `yaml:"map_size"`
	TurnTimer int    `yaml:"turn_timer"`
	TurnMode  string `yaml:"turn_mode"`
}

type WinCondConfig struct {
	DominanceTurnsRequired int `yaml:"dominance_turns_required"`
}

// LoadBalance reads and parses the balance YAML file.
func LoadBalance(path string) (*BalanceData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: failed to read balance file %s: %w", path, err)
	}

	var balance BalanceData
	if err := yaml.Unmarshal(data, &balance); err != nil {
		return nil, fmt.Errorf("config: failed to parse balance file %s: %w", path, err)
	}

	return &balance, nil
}
