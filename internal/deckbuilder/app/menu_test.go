package deckbuilder

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	releasenotes "github.com/HybridUofA/casters-compendium/docs/releases"
)

func TestBuildMainMenuDisplaysEnabledSimulatorPrototype(t *testing.T) {
	window := test.NewWindow(nil)
	defer window.Close()

	menu := buildMainMenu(window, mainMenuActions{PlayGame: func() {}})
	button := findButtonText(menu, "Play a Game (Hotseat)")
	if button == nil {
		t.Fatal("main menu does not display the simulator hotseat action")
	}
	if button.Disabled() {
		t.Fatal("simulator prototype is disabled despite having an action")
	}
}

func TestBuildMainMenuDisplaysSoloVsTsubasa(t *testing.T) {
	window := test.NewWindow(nil)
	defer window.Close()

	menu := buildMainMenu(window, mainMenuActions{PlaySoloVsTsubasa: func() {}})
	button := findButtonText(menu, "Vs Tsubasa")
	if button == nil {
		t.Fatal("main menu does not display Vs Tsubasa")
	}
	if button.Disabled() {
		t.Fatal("Vs Tsubasa is disabled despite having an action")
	}
}

func TestBuildMainMenuDisplaysOnlineAndDeckButtons(t *testing.T) {
	window := test.NewWindow(nil)
	defer window.Close()

	menu := buildMainMenu(window, mainMenuActions{})
	for _, label := range []string{"Online", "Decks"} {
		if findButtonText(menu, label) == nil {
			t.Fatalf("main menu does not display the %q button", label)
		}
	}
}

func TestBuildMainMenuDisplaysUtilityButtons(t *testing.T) {
	window := test.NewWindow(nil)
	defer window.Close()

	menu := buildMainMenu(window, mainMenuActions{})
	for _, label := range []string{"Check for Updates", "How to Use", "Settings", "Quit"} {
		if findButtonText(menu, label) == nil {
			t.Fatalf("main menu does not display the %q button", label)
		}
	}
}

func TestOnlineMenuListsOnlinePlayActions(t *testing.T) {
	var lobby, host, join bool
	menu := newOnlineMenu(mainMenuActions{
		OnlineLobby:    func() { lobby = true },
		HostOnlineGame: func() { host = true },
		JoinOnlineGame: func() { join = true },
	})

	for _, label := range []string{"Online Lobby", "Host Online Game", "Join Online Game"} {
		item := findMenuItem(menu, label)
		if item == nil {
			t.Fatalf("online menu does not list %q", label)
		}
		item.Action()
	}
	if !lobby || !host || !join {
		t.Fatal("online menu items did not run their actions")
	}
}

func TestDeckMenuListsDeckActions(t *testing.T) {
	var editor bool
	menu := newDeckMenu(mainMenuActions{OpenDeckEditor: func() { editor = true }})

	for _, label := range []string{
		"Open Deck Editor",
		"Make a New Deck",
		"Load a Deck",
		"Generate Deck Image from Decklist",
		"Generate Decklist File",
	} {
		if findMenuItem(menu, label) == nil {
			t.Fatalf("deck menu does not list %q", label)
		}
	}
	findMenuItem(menu, "Open Deck Editor").Action()
	if !editor {
		t.Fatal("Open Deck Editor menu item did not run its action")
	}
}

func TestBuildMainMenuDisplaysLocalChangelog(t *testing.T) {
	notes, err := releasenotes.All()
	if err != nil || len(notes) == 0 {
		t.Fatalf("bundled release notes unavailable: %v", err)
	}

	window := test.NewWindow(nil)
	defer window.Close()

	// The changelog is shown inline; its version picker starts on the newest release.
	menu := buildMainMenu(window, mainMenuActions{})
	if findButtonText(menu, notes[0].Version) == nil {
		t.Fatalf("main menu does not display the changelog for %q", notes[0].Version)
	}
}

func findMenuItem(menu *fyne.Menu, label string) *fyne.MenuItem {
	for _, item := range menu.Items {
		if item.Label == label {
			return item
		}
	}
	return nil
}

func findButtonText(object fyne.CanvasObject, text string) *widget.Button {
	if button, ok := object.(*widget.Button); ok && button.Text == text {
		return button
	}
	if container, ok := object.(*fyne.Container); ok {
		for _, child := range container.Objects {
			if button := findButtonText(child, text); button != nil {
				return button
			}
		}
	}
	return nil
}
