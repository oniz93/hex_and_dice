package game

import (
	"github.com/teomiscia/hexbattle/internal/dice"
	"github.com/teomiscia/hexbattle/internal/model"
	"github.com/teomiscia/hexbattle/internal/ws"
)

// ResolveTroopCombat resolves a full combat exchange between an attacker troop and a defender troop.
// Returns the combat result delta and any troop destroyed deltas.
func ResolveTroopCombat(gs *GameState, roller *dice.Roller, attacker, defender *model.Troop) (*ws.CombatResultData, []*ws.TroopDestroyedData) {
	var destroyed []*ws.TroopDestroyedData

	// --- Primary attack ---
	attackerTerrain := gs.GetTerrainAt(attacker.Hex)
	defenderTerrain := gs.GetTerrainAt(defender.Hex)

	naturalRoll := roller.D20()
	atkModifier := attacker.ATK + model.GetTerrainInfo(attackerTerrain).ATKModifier
	totalRoll := naturalRoll + atkModifier
	targetDEF := defender.DEF + model.GetTerrainInfo(defenderTerrain).DEFModifier

	isCrit := naturalRoll == 20
	isFumble := naturalRoll == 1

	var hit bool
	if isCrit {
		hit = true
	} else if isFumble {
		hit = false
	} else {
		hit = totalRoll >= targetDEF
	}

	damageDealt := 0
	if hit {
		dn, err := TroopDamageDice(attacker)
		if err == nil {
			if isCrit {
				damageDealt = roller.RollDamage(dn) * 2
			} else {
				damageDealt = roller.RollDamage(dn)
			}
		}
		defender.TakeDamage(damageDealt)
	}

	defenderKill := !defender.IsAlive()
	if defenderKill {
		destroyed = append(destroyed, &ws.TroopDestroyedData{
			UnitID: defender.ID,
			HexQ:   defender.Hex.Q,
			HexR:   defender.Hex.R,
			HexS:   defender.Hex.S,
			Cause:  "combat",
		})
	}

	result := &ws.CombatResultData{
		AttackerID:  attacker.ID,
		DefenderID:  defender.ID,
		HitRoll:     totalRoll,
		NaturalRoll: naturalRoll,
		Hit:         hit,
		DamageRoll:  damageDealt,
		Damage:      damageDealt,
		DefenderHP:  defender.CurrentHP,
		Killed:      defenderKill,
		Crit:        isCrit,
		Fumble:      isFumble,
		AttackerHP:  attacker.CurrentHP,
	}

	// --- Counterattack ---
	// Triggered when:
	// 1. Melee attacker (range=1) attacks a melee defender (range=1), OR
	// 2. Attacker rolled a fumble (natural 1) — defender gets free counter regardless
	// Only if defender is still alive
	shouldCounter := false
	if defender.IsAlive() {
		if isFumble {
			shouldCounter = true
		} else if attacker.IsMelee() && defender.IsMelee() {
			shouldCounter = true
		}
	}

	if shouldCounter {
		result.HasCounter = true

		counterNatural := roller.D20()
		counterAtkMod := defender.ATK + model.GetTerrainInfo(defenderTerrain).ATKModifier
		counterTotal := counterNatural + counterAtkMod
		counterTargetDEF := attacker.DEF + model.GetTerrainInfo(attackerTerrain).DEFModifier

		counterCrit := counterNatural == 20
		counterFumble := counterNatural == 1

		var counterHit bool
		if counterCrit {
			counterHit = true
		} else if counterFumble {
			counterHit = false
		} else {
			counterHit = counterTotal >= counterTargetDEF
		}

		counterDamage := 0
		if counterHit {
			dn, err := TroopDamageDice(defender)
			if err == nil {
				counterDamage = roller.RollHalfDamage(dn)
			}
			attacker.TakeDamage(counterDamage)
		}

		result.CounterHitRoll = counterTotal
		result.CounterNatural = counterNatural
		result.CounterHit = counterHit
		result.CounterDamage = counterDamage
		result.AttackerHP = attacker.CurrentHP
		result.AttackerKilled = !attacker.IsAlive()

		if result.AttackerKilled {
			destroyed = append(destroyed, &ws.TroopDestroyedData{
				UnitID: attacker.ID,
				HexQ:   attacker.Hex.Q,
				HexR:   attacker.Hex.R,
				HexS:   attacker.Hex.S,
				Cause:  "combat",
			})
		}
	}

	// Mark both as having been in combat
	attacker.WasInCombat = true
	defender.WasInCombat = true
	attacker.HasAttacked = true

	return result, destroyed
}

// ResolveStructureAttack resolves a troop attacking a structure.
// Returns the structure attacked delta and optionally a troop destroyed delta (from fumble counter).
func ResolveStructureAttack(gs *GameState, roller *dice.Roller, attacker *model.Troop, structure *model.Structure) (*ws.StructureAttackedData, []*ws.TroopDestroyedData) {
	var destroyed []*ws.TroopDestroyedData

	attackerTerrain := gs.GetTerrainAt(attacker.Hex)
	naturalRoll := roller.D20()
	atkModifier := attacker.ATK + model.GetTerrainInfo(attackerTerrain).ATKModifier
	totalRoll := naturalRoll + atkModifier
	targetDEF := structure.DEF

	isCrit := naturalRoll == 20
	isFumble := naturalRoll == 1

	var hit bool
	if isCrit {
		hit = true
	} else if isFumble {
		hit = false
	} else {
		hit = totalRoll >= targetDEF
	}

	damageDealt := 0
	if hit {
		dn, err := TroopDamageDice(attacker)
		if err == nil {
			if isCrit {
				damageDealt = roller.RollDamage(dn) * 2
			} else {
				damageDealt = roller.RollDamage(dn)
			}
			// Anti-structure multiplier (Mech does 2x vs structures)
			damageDealt *= AntiStructureMultiplier(attacker.Type)
		}
	}

	captured := false
	newOwner := ""
	if hit {
		if structure.TakeDamage(damageDealt) {
			// Structure reaches 0 HP: capture it
			structure.Capture(attacker.OwnerID)
			captured = true
			newOwner = attacker.OwnerID
		}
	}

	attacker.HasAttacked = true
	attacker.WasInCombat = true

	result := &ws.StructureAttackedData{
		StructureID: structure.ID,
		AttackerID:  attacker.ID,
		HitRoll:     totalRoll,
		Damage:      damageDealt,
		StructureHP: structure.CurrentHP,
		Captured:    captured,
		NewOwner:    newOwner,
	}

	return result, destroyed
}

// ResolveStructureFire resolves a structure attacking a troop (auto-attack phase).
func ResolveStructureFire(gs *GameState, roller *dice.Roller, structure *model.Structure, target *model.Troop) *ws.StructureFiresData {
	naturalRoll := roller.D20()
	atkModifier := structure.ATK
	totalRoll := naturalRoll + atkModifier

	targetTerrain := gs.GetTerrainAt(target.Hex)
	targetDEF := target.DEF + model.GetTerrainInfo(targetTerrain).DEFModifier

	// Neutral structures have reduced ATK
	if structure.IsNeutral() && Balance != nil {
		atkModifier -= Balance.NeutralMod.ATKReduction
		totalRoll = naturalRoll + atkModifier
	}

	isCrit := naturalRoll == 20
	isFumble := naturalRoll == 1

	var hit bool
	if isCrit {
		hit = true
	} else if isFumble {
		hit = false
	} else {
		hit = totalRoll >= targetDEF
	}

	damageDealt := 0
	if hit {
		dn, err := StructureDamageDice(structure)
		if err == nil {
			// Neutral structures have reduced damage dice
			if structure.IsNeutral() && Balance != nil {
				for i := 0; i < Balance.NeutralMod.DamageStepDown; i++ {
					dn = dn.StepDown()
				}
			}
			if isCrit {
				damageDealt = roller.RollDamage(dn) * 2
			} else {
				damageDealt = roller.RollDamage(dn)
			}
		}
		target.TakeDamage(damageDealt)
	}

	return &ws.StructureFiresData{
		StructureID: structure.ID,
		TargetID:    target.ID,
		HitRoll:     totalRoll,
		Damage:      damageDealt,
		TargetHP:    target.CurrentHP,
		Killed:      !target.IsAlive(),
	}
}

// FindStructureTarget finds the best target for a structure to auto-attack.
// Chooses the closest enemy troop in range (random tiebreak via roller).
func FindStructureTarget(gs *GameState, roller *dice.Roller, structure *model.Structure) *model.Troop {
	var candidates []*model.Troop

	for _, troop := range gs.Troops {
		if !troop.IsAlive() {
			continue
		}

		// Player-owned structures attack enemy troops only.
		// Neutral structures attack troops of ANY player.
		if !structure.IsNeutral() && troop.OwnerID == structure.OwnerID {
			continue
		}

		dist := troop.Hex.Distance(structure.Hex)
		if dist <= structure.Range {
			candidates = append(candidates, troop)
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	// Find minimum distance
	minDist := candidates[0].Hex.Distance(structure.Hex)
	for _, c := range candidates[1:] {
		d := c.Hex.Distance(structure.Hex)
		if d < minDist {
			minDist = d
		}
	}

	// Filter to closest only
	var closest []*model.Troop
	for _, c := range candidates {
		if c.Hex.Distance(structure.Hex) == minDist {
			closest = append(closest, c)
		}
	}

	// Random tiebreak
	if len(closest) == 1 {
		return closest[0]
	}
	return closest[roller.Roll(len(closest))-1]
}
