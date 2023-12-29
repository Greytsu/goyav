package track

import "github.com/zmb3/spotify/v2"

func MapSpotifyTrackToGoyavTrack(track spotify.FullTrack) Track {
	return Track{
		Name:           track.Name,
		SpotifyId:      track.ID,
		SpotifyAlbumId: track.Album.ID,
		SpotifyUrl:     track.ExternalURLs["spotify"],
		ImageBigUrl:    track.Album.Images[0].URL,
		ImageMediumUrl: track.Album.Images[1].URL,
		ImageSmallUrl:  track.Album.Images[2].URL,
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
		Name:           track.Name,
		SpotifyId:      track.ID,
		SpotifyAlbumId: track.Album.ID,
		SpotifyUrl:     track.ExternalURLs["spotify"],
		ImageBigUrl:    track.Album.Images[0].URL,
		ImageMediumUrl: track.Album.Images[1].URL,
		ImageSmallUrl:  track.Album.Images[2].URL,
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
		seeds = append(seeds, track.SpotifyId)
	}
	return seeds
}
