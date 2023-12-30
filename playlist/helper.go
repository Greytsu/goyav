package playlist

import "errors"

var UnknownAttributeError = errors.New("unknown AttributeType")

// IsValidAttributeValue checks if the AttributeValue is within the valid bounds for the given AttributeType
func IsValidAttributeValue(attrType AttributeType, attrValue float32) (bool, error) {
	switch attrType {
	case Acousticness:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	case Danceability:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	case Energy:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	case Instrumentalness:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	case Key:
		return attrValue >= 0.0 && attrValue <= 11.0, nil
	case Liveness:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	case Loudness:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	case Mode:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	case Popularity:
		return attrValue >= 0.0 && attrValue <= 100.0, nil
	case Speechiness:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	case Tempo:
		return attrValue >= 0.0 && attrValue <= 300.0, nil
	case Valence:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	default:
		return false, UnknownAttributeError
	}
}
