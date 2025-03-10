package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"slices"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var instance *AssetManager

// AssetManager loads and stores assets from files.
type AssetManager struct {
	textures map[string]rl.Texture2D
	sounds   map[string]rl.Sound
	music    map[string]rl.Music
	fs       embed.FS
}

func GetAssetManagerInstance() *AssetManager {
	if instance == nil {
		instance = &AssetManager{
			textures: map[string]rl.Texture2D{},
			sounds:   map[string]rl.Sound{},
			music:    map[string]rl.Music{},
		}
	}

	return instance
}

func (am *AssetManager) LoadAssets(embedFS embed.FS) {
	soundExts := []string{".mp3", ".ogg", ".wav"}
	texExts := []string{".jpg", ".jpeg", ".png"}

	am.UnloadTextures()
	am.textures = make(map[string]rl.Texture2D)
	am.sounds = make(map[string]rl.Sound)
	am.fs = embedFS

	// Walk through all subdirectories and load all assets
	err := fs.WalkDir(embedFS, "assets", func(path string, entry fs.DirEntry, err error) error {
		log.Println("Asset manager processing path: ", path)

		if err != nil {
			log.Printf("Error reading asset: %v\n", err)
			return nil
		}

		if entry.IsDir() {
			return nil // Skip directories
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == "" {
			log.Printf("Cannot determine file format of %s. Skipping\n", path)
			return nil
		}

		if slices.Contains(texExts, ext) {
			if texture, err := am.loadTextureFromEmbed(path, ext); err == nil {
				am.textures[path] = texture
			} else {
				log.Printf("Skipping non-image asset: %s, error: %v\n", path, err)
			}
			return nil
		}

		// TODO: Figure out why loading the music as rl.Music causes a crash for now load them as SFX
		isMusic := strings.Index(path, "assets/music") > -1
		isSfx := strings.Index(path, "assets/sfx") > -1 || true

		if isSfx && slices.Contains(soundExts, ext) {
			if sound, err := am.loadSoundFromEmbed(path, ext); err == nil {
				am.sounds[path] = sound
			} else {
				log.Printf("Skipping non-sound asset: %s, error: %v\n", path, err)
			}
			return nil
		}

		if isMusic && slices.Contains(soundExts, ext) {
			if music, err := am.loadMusicFromEmbed(path, ext); err == nil {
				am.music[path] = music
			} else {
				log.Printf("Skipping non-music asset: %s, error: %v\n", path, err)
			}
			return nil
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to load assets: %v\n", err)
	}

	log.Printf("Loaded assets!\n%+v\n", am.textures)
}

// loadTexture reads an embedded file and loads it as a Raylib texture.
func (am *AssetManager) loadTextureFromEmbed(path string, fileType string) (rl.Texture2D, error) {
	log.Println("Asset manager attempting to load file: ", path, " with type ", fileType)
	data, err := am.fs.ReadFile(path)
	if err != nil {
		return rl.Texture2D{}, fmt.Errorf("failed to read embedded asset: %v", err)
	}

	raylibImage := rl.LoadImageFromMemory(fileType, data, int32(len(data)))
	if raylibImage.Width == 0 || raylibImage.Height == 0 {
		return rl.Texture2D{}, fmt.Errorf("failed to load image: %s", path)
	}

	// Convert the Raylib image into a texture
	texture := rl.LoadTextureFromImage(raylibImage)
	rl.UnloadImage(raylibImage) // Cleanup Raylib Image after loading texture

	return texture, nil
}

func (am *AssetManager) loadSoundFromEmbed(path string, fileType string) (rl.Sound, error) {
	data, err := am.fs.ReadFile(path)
	if err != nil {
		return rl.Sound{}, fmt.Errorf("failed to read embedded asset: %s", err)
	}

	raylibWave := rl.LoadWaveFromMemory(fileType, data, int32(len(data)))
	if raylibWave.FrameCount == 0 {
		return rl.Sound{}, fmt.Errorf("failed to load wave from path: %s", path)
	}

	//TODO: There might be an error here will need to revisit when we try to get audio running
	sound := rl.LoadSoundFromWave(raylibWave)
	rl.UnloadWave(raylibWave)

	return sound, nil
}

func (am *AssetManager) loadMusicFromEmbed(path string, fileType string) (rl.Music, error) {
	data, err := am.fs.ReadFile(path)
	if err != nil {
		return rl.Music{}, fmt.Errorf("Failed to read embedded asset: %s", err)
	}

	stream := rl.LoadMusicStreamFromMemory(fileType, data, int32(len(data)))
	if stream.Stream.SampleRate <= 0 {
		return rl.Music{}, fmt.Errorf("failed to load music stream from path: %s", path)
	}

	return stream, nil
}

// GetTexture retrieves a texture by file path.
func (am *AssetManager) GetTexture(path string) (rl.Texture2D, bool) {
	texture, found := am.textures[path]
	return texture, found
}

func (am *AssetManager) GetSound(path string) (rl.Sound, bool) {
	sound, found := am.sounds[path]
	return sound, found
}

func (am *AssetManager) GetMusic(path string) (*rl.Music, bool) {
	music, found := am.music[path]
	return &music, found
}

// UnloadTextures releases all textures.
func (am *AssetManager) UnloadTextures() {
	hasSounds := am.sounds != nil
	hasTextures := am.textures != nil

	if !hasSounds && !hasTextures {
		return
	}

	if hasTextures {
		for _, texture := range am.textures {
			rl.UnloadTexture(texture)
		}
	}

	if hasSounds {
		for _, sound := range am.sounds {
			rl.UnloadSound(sound)
		}
	}

	am.textures = nil
	am.sounds = nil
}
