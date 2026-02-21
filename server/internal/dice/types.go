package dice

// Result captures the full outcome of a combat hit roll.
type Result struct {
	NaturalRoll int  `json:"natural_roll"` // The raw D20 result (1-20)
	TotalRoll   int  `json:"total_roll"`   // D20 + ATK modifier + terrain modifier
	TargetDEF   int  `json:"target_def"`   // DEF score that needed to be met
	Hit         bool `json:"hit"`          // Whether the attack connected
	IsCrit      bool `json:"is_crit"`      // Natural 20
	IsFumble    bool `json:"is_fumble"`    // Natural 1
}

// DamageResult captures the outcome of a damage roll.
type DamageResult struct {
	Notation string `json:"notation"` // e.g. "2D6+2"
	Total    int    `json:"total"`    // Total damage dealt
	IsHalf   bool   `json:"is_half"`  // True if this was a half-damage counterattack roll
}

// CombatResult captures the full outcome of a combat exchange.
type CombatResult struct {
	AttackerID string `json:"attacker_id"`
	DefenderID string `json:"defender_id"`

	// Primary attack
	HitResult    Result       `json:"hit_result"`
	DamageResult DamageResult `json:"damage_result"`
	DefenderHP   int          `json:"defender_hp"`
	DefenderKill bool         `json:"defender_kill"`

	// Counterattack (may be nil if no counter triggered)
	HasCounter       bool         `json:"has_counter"`
	CounterHitResult Result       `json:"counter_hit_result,omitempty"`
	CounterDamage    DamageResult `json:"counter_damage,omitempty"`
	AttackerHP       int          `json:"attacker_hp"`
	AttackerKill     bool         `json:"attacker_kill"`
}
