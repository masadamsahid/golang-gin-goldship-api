package googlemap

import (
	"context"
	"os"

	"googlemaps.github.io/maps"
)

var GOOGLE_MAP_API_KEY string

func InitGoogleMapAPI() {
	GOOGLE_MAP_API_KEY = os.Getenv("GOOGLE_MAP_API_KEY")
}

func CalculateDistance(distanceMatrixRequest *maps.DistanceMatrixRequest) (*maps.Distance, error) {

	mapAPIClient, err := maps.NewClient(maps.WithAPIKey(GOOGLE_MAP_API_KEY))
	if err != nil {
		return nil, err
	}

	resp, err := mapAPIClient.DistanceMatrix(context.Background(), distanceMatrixRequest)
	if err != nil {
		return nil, err
	}

	// str, err := json.Marshal(resp)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Println("AAAAAAAA:", string(str))

	// for i, row := range resp.Rows {
	// 	origin := resp.OriginAddresses[i]
	// 	for j, element := range row.Elements {
	// 		destination := resp.DestinationAddresses[j]
	// 		if element.Status == "OK" {
	// 			fmt.Printf("Distance from %s to %s is %+v and the duration is %+v\n",
	// 				origin,
	// 				destination,
	// 				element.Distance,
	// 				element.Duration,
	// 			)
	// 		} else {
	// 			fmt.Printf("Could not calculate distance from %s to %s. Status: %s\n",
	// 				origin,
	// 				destination,
	// 				element.Status,
	// 			)
	// 		}
	// 	}
	// }

	return &resp.Rows[0].Elements[0].Distance, nil

}
