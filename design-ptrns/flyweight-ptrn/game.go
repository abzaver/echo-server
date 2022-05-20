// Client code
package main

type game struct {
	terrorists        []*Player
	counterTerrorists []*Player
}

func newGame() *game {
	return &game{
		terrorists:        make([]*Player, 1),
		counterTerrorists: make([]*Player, 1),
	}
}

func (c *game) addTerrorist(dressType string) {
	player := newPlayer("T", dressType)
	player.newLocation(0, 0)
	c.terrorists = append(c.terrorists, player)
}

func (c *game) addCounterTerrorist(dressType string) {
	player := newPlayer("CT", dressType)
	player.newLocation(100, 100)
	c.counterTerrorists = append(c.counterTerrorists, player)
}
