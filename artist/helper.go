package artist

import "github.com/zmb3/spotify/v2"

func MapSpotifyArtistToGoyavArtist(spotifyArtist spotify.FullArtist) Artist {
	return Artist{
		Name:           spotifyArtist.Name,
		SpotifyId:      spotifyArtist.ID,
		SpotifyUrl:     spotifyArtist.ExternalURLs["spotify"],
		Genres:         spotifyArtist.Genres,
		ImageBigUrl:    spotifyArtist.Images[0].URL,
		ImageMediumUrl: spotifyArtist.Images[1].URL,
		ImageSmallUrl:  spotifyArtist.Images[2].URL,
	}
}

func MapSpotifyArtistsToGoyavArtists(spotifyArtists []spotify.FullArtist) []Artist {
	var artists []Artist

	for _, spotifyArtist := range spotifyArtists {
		artists = append(artists, MapSpotifyArtistToGoyavArtist(spotifyArtist))
	}

	return artists
}

func MapToTrackSeeds(artists []Artist) []spotify.ID {
	var seeds []spotify.ID
	for _, artist := range artists {
		seeds = append(seeds, artist.SpotifyId)
	}
	return seeds
}

func ExtractGenresFromArtists(artists []Artist) []string {
	var genres []string
	for _, artist := range artists {
		if len(artist.Genres) > 0 {
			genres = append(genres, artist.Genres[0])
		}
	}
	return genres
}
