package osint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TwiN/go-color"
	"github.com/manifoldco/promptui"
)

const favoritesFile = "favorites.json"

// FavoriteSatellite represents a saved favorite satellite.
type FavoriteSatellite struct {
	SatelliteName string `json:"satellite_name"`
	NORADID       string `json:"norad_id"`
	Country       string `json:"country,omitempty"`
	ObjectType    string `json:"object_type,omitempty"`
	AddedDate    string `json:"added_date"`
}

// FavoritesList represents the collection of favorite satellites.
type FavoritesList struct {
	Favorites []FavoriteSatellite `json:"favorites"`
}

// getFavoritesPath returns the full path to the favorites file.
func getFavoritesPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return favoritesFile
	}
	favoritesDir := filepath.Join(homeDir, ".satintel")
	os.MkdirAll(favoritesDir, 0755)
	return filepath.Join(favoritesDir, favoritesFile)
}

// LoadFavorites reads the favorites list from the JSON file.
func LoadFavorites() ([]FavoriteSatellite, error) {
	favoritesPath := getFavoritesPath()
	
	data, err := os.ReadFile(favoritesPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, return empty list
			return []FavoriteSatellite{}, nil
		}
		return nil, fmt.Errorf("failed to read favorites file: %w", err)
	}

	var favoritesList FavoritesList
	if err := json.Unmarshal(data, &favoritesList); err != nil {
		return nil, fmt.Errorf("failed to parse favorites file: %w", err)
	}

	return favoritesList.Favorites, nil
}

// SaveFavorites writes the favorites list to the JSON file.
func SaveFavorites(favorites []FavoriteSatellite) error {
	favoritesPath := getFavoritesPath()
	
	favoritesList := FavoritesList{
		Favorites: favorites,
	}

	data, err := json.MarshalIndent(favoritesList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal favorites: %w", err)
	}

	if err := os.WriteFile(favoritesPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write favorites file: %w", err)
	}

	return nil
}

// AddFavorite adds a satellite to the favorites list.
func AddFavorite(satName, noradID, country, objectType string) error {
	favorites, err := LoadFavorites()
	if err != nil {
		return err
	}

	// Check if already exists
	for _, fav := range favorites {
		if fav.NORADID == noradID {
			return fmt.Errorf("satellite with NORAD ID %s is already in favorites", noradID)
		}
	}

	newFavorite := FavoriteSatellite{
		SatelliteName: satName,
		NORADID:       noradID,
		Country:       country,
		ObjectType:    objectType,
		AddedDate:     time.Now().Format("2006-01-02 15:04:05"),
	}

	favorites = append(favorites, newFavorite)
	return SaveFavorites(favorites)
}

// RemoveFavorite removes a satellite from the favorites list by NORAD ID.
func RemoveFavorite(noradID string) error {
	favorites, err := LoadFavorites()
	if err != nil {
		return err
	}

	var updatedFavorites []FavoriteSatellite
	found := false
	for _, fav := range favorites {
		if fav.NORADID != noradID {
			updatedFavorites = append(updatedFavorites, fav)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("satellite with NORAD ID %s not found in favorites", noradID)
	}

	return SaveFavorites(updatedFavorites)
}

// IsFavorite checks if a satellite is in the favorites list.
func IsFavorite(noradID string) (bool, error) {
	favorites, err := LoadFavorites()
	if err != nil {
		return false, err
	}

	for _, fav := range favorites {
		if fav.NORADID == noradID {
			return true, nil
		}
	}

	return false, nil
}

// SelectFromFavorites displays a menu to select from saved favorites.
func SelectFromFavorites() string {
	favorites, err := LoadFavorites()
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to load favorites: "+err.Error()))
		return ""
	}

	if len(favorites) == 0 {
		fmt.Println(color.Ize(color.Yellow, "  [!] No favorites saved yet"))
		fmt.Println(color.Ize(color.Cyan, "  [*] Add favorites by selecting 'Save to Favorites' after choosing a satellite"))
		return ""
	}

	var menuItems []string
	for _, fav := range favorites {
		info := fmt.Sprintf("%s (%s)", fav.SatelliteName, fav.NORADID)
		if fav.Country != "" {
			info += fmt.Sprintf(" - %s", fav.Country)
		}
		if fav.ObjectType != "" {
			info += fmt.Sprintf(" [%s]", fav.ObjectType)
		}
		info += fmt.Sprintf(" (Added: %s)", fav.AddedDate)
		menuItems = append(menuItems, info)
	}

	menuItems = append(menuItems, "❌ Cancel")

	prompt := promptui.Select{
		Label: fmt.Sprintf("Select from Favorites ⭐ (%d saved)", len(favorites)),
		Items: menuItems,
		Size:  15,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] PROMPT FAILED"))
		return ""
	}

	if idx >= len(favorites) {
		// Cancel
		return ""
	}

	selected := favorites[idx]
	return fmt.Sprintf("%s (%s)", selected.SatelliteName, selected.NORADID)
}

// ManageFavorites provides an interactive menu to manage favorites.
func ManageFavorites() {
	favorites, err := LoadFavorites()
	if err != nil {
		fmt.Println(color.Ize(color.Red, "  [!] ERROR: Failed to load favorites: "+err.Error()))
		return
	}

	if len(favorites) == 0 {
		fmt.Println(color.Ize(color.Yellow, "  [!] No favorites saved yet"))
		return
	}

	menuItems := []string{
		"View All Favorites",
		"Remove Favorite",
		"Clear All Favorites",
		"Back",
	}

	prompt := promptui.Select{
		Label: fmt.Sprintf("Manage Favorites ⭐ (%d saved)", len(favorites)),
		Items: menuItems,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return
	}

	switch idx {
	case 0: // View All Favorites
		fmt.Println(color.Ize(color.Cyan, "\n  Your Favorites:"))
		fmt.Println(strings.Repeat("-", 70))
		for i, fav := range favorites {
			fmt.Printf("%d. %s (%s)\n", i+1, fav.SatelliteName, fav.NORADID)
			if fav.Country != "" {
				fmt.Printf("   Country: %s\n", fav.Country)
			}
			if fav.ObjectType != "" {
				fmt.Printf("   Type: %s\n", fav.ObjectType)
			}
			fmt.Printf("   Added: %s\n", fav.AddedDate)
			if i < len(favorites)-1 {
				fmt.Println()
			}
		}
		fmt.Println(strings.Repeat("-", 70))

	case 1: // Remove Favorite
		var removeItems []string
		for _, fav := range favorites {
			removeItems = append(removeItems, fmt.Sprintf("%s (%s)", fav.SatelliteName, fav.NORADID))
		}
		removeItems = append(removeItems, "Cancel")

		removePrompt := promptui.Select{
			Label: "Select Favorite to Remove",
			Items: removeItems,
		}

		removeIdx, _, err := removePrompt.Run()
		if err != nil || removeIdx >= len(favorites) {
			return
		}

		selected := favorites[removeIdx]
		if err := RemoveFavorite(selected.NORADID); err != nil {
			fmt.Println(color.Ize(color.Red, "  [!] ERROR: "+err.Error()))
		} else {
			fmt.Println(color.Ize(color.Green, fmt.Sprintf("  [+] Removed %s from favorites", selected.SatelliteName)))
		}

	case 2: // Clear All Favorites
		confirmPrompt := promptui.Prompt{
			Label:     "Are you sure you want to clear all favorites? (yes/no)",
			Default:   "no",
			AllowEdit: true,
		}

		confirm, err := confirmPrompt.Run()
		if err != nil {
			return
		}

		if strings.ToLower(strings.TrimSpace(confirm)) == "yes" {
			if err := SaveFavorites([]FavoriteSatellite{}); err != nil {
				fmt.Println(color.Ize(color.Red, "  [!] ERROR: "+err.Error()))
			} else {
				fmt.Println(color.Ize(color.Green, "  [+] All favorites cleared"))
			}
		}
	}
}

