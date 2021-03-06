package commands

type steamGames struct {
	Applist struct {
		Apps []struct {
			Appid int    `json:"appid"`
			Name  string `json:"name"`
		} `json:"apps"`
	} `json:"applist"`
}

type affirmation struct {
	Affirmation string `json:"affirmation"`
}

type kanye struct {
	Quote string `json:"quote"`
}

type advice struct {
	Slip struct {
		ID     int    `json:"id"`
		Advice string `json:"advice"`
	} `json:"slip"`
}

type doggo []struct {
	Breeds []struct {
		Weight struct {
			Imperial string `json:"imperial"`
			Metric   string `json:"metric"`
		} `json:"weight"`
		Height struct {
			Imperial string `json:"imperial"`
			Metric   string `json:"metric"`
		} `json:"height"`
		ID               int    `json:"id"`
		Name             string `json:"name"`
		BredFor          string `json:"bred_for"`
		BreedGroup       string `json:"breed_group"`
		LifeSpan         string `json:"life_span"`
		Temperament      string `json:"temperament"`
		Origin           string `json:"origin"`
		ReferenceImageID string `json:"reference_image_id"`
	} `json:"breeds"`
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type joke struct {
	ID   string `json:"id"`
	Joke string `json:"joke"`
}

type pickupLine struct {
	Category string `json:"category"`
	Joke     string `json:"joke"`
}

type wyr struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type wtp struct {
	Data struct {
		Type      []string `json:"Type"`
		Abilities []string `json:"abilities"`
		ASCII     string   `json:"ascii"`
		Height    float64  `json:"height"`
		ID        int      `json:"id"`
		Link      string   `json:"link"`
		Name      string   `json:"name"`
		Weight    int      `json:"weight"`
	} `json:"Data"`
	Answer   string `json:"answer"`
	Question string `json:"question"`
}

type albumPicker struct {
	Genres      string `json:"Genres"`
	Descriptors string `json:"Descriptors"`
	Artist      string `json:"Artist"`
	AlbumName   string `json:"Album_Name"`
	SecGenres   string `json:"Sec_Genres"`
	URL         string `json:"url"`
	//SpotifyAlbumURL string `json:"spotify_album_url"`
}
