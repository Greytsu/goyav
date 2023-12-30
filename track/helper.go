package track

import "github.com/zmb3/spotify/v2"

func MapSpotifyTrackToGoyavTrack(track spotify.FullTrack) Track {
	return Track{
		TrackName:        track.Name,
		TrackId:          track.ID,
		SpotifyAlbumName: track.Album.Name,
		SpotifyAlbumId:   track.Album.ID,
		SpotifyUrl:       track.ExternalURLs["spotify"],
		ImageUrls: ImageUrls{
			Big:    track.Album.Images[0].URL,
			Medium: track.Album.Images[1].URL,
			Small:  track.Album.Images[2].URL,
		},
	}
}

func MapSpotifyTracksToGoyavTracks(spotifyTracks []spotify.FullTrack) []Track {
	var tracks []Track

	for _, spotifyArtist := range spotifyTracks {
		tracks = append(tracks, MapSpotifyTrackToGoyavTrack(spotifyArtist))
	}

	return tracks
}

func MapSpotifySimpleTrackToGoyavTrack(track spotify.SimpleTrack) Track {
	return Track{
		TrackName:        track.Name,
		TrackId:          track.ID,
		ArtistName:       track.Artists[0].Name,
		ArtistId:         track.Artists[0].ID,
		SpotifyAlbumName: track.Album.Name,
		SpotifyAlbumId:   track.Album.ID,
		SpotifyUrl:       track.ExternalURLs["spotify"],
		ImageUrls: ImageUrls{
			Big:    track.Album.Images[0].URL,
			Medium: track.Album.Images[1].URL,
			Small:  track.Album.Images[2].URL,
		},
		Contributors: []string{},
	}
}

func MapSpotifySimpleTracksToGoyavTracks(spotifyTracks []spotify.SimpleTrack) []Track {
	var tracks []Track

	for _, spotifyArtist := range spotifyTracks {
		tracks = append(tracks, MapSpotifySimpleTrackToGoyavTrack(spotifyArtist))
	}

	return tracks
}

func MapToTrackSeeds(tracks []Track) []spotify.ID {
	var seeds []spotify.ID
	for _, track := range tracks {
		seeds = append(seeds, track.TrackId)
	}
	return seeds
}
