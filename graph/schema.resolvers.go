package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/VarunAttarde22/hackernews/graph/generated"
	"github.com/VarunAttarde22/hackernews/graph/model"
)

func (r *mutationResolver) CreateNode(ctx context.Context, input model.NewNode) (*model.Node, error) {
	var inputNetworks []*model.Ips
	var inputRoles []*model.Roles
	for _, val := range input.Networks {
		ip := &model.Ips{
			IPType:  val.IPType,
			IP:      val.IP,
			Netmask: val.Netmask,
			Gateway: val.Gateway,
		}
		inputNetworks = append(inputNetworks, ip)
	}
	for _, val := range input.Roles {
		role := &model.Roles{
			RoleType: val.RoleType,
			Label:    val.Label,
		}
		inputRoles = append(inputRoles, role)
	}
	// fmt.Println(inputNetworks)
	node := &model.Node{
		IP:          input.IP,
		Serial:      input.Serial,
		Model:       input.Model,
		Credentials: input.Credentials,
		Template:    input.Template,
		Networks:    inputNetworks,
		Roles:       inputRoles,
	}
	r.nodes = append(r.nodes, (*model.Node)(node))
	return node, nil
}

func (r *queryResolver) Nodes(ctx context.Context) ([]*model.Node, error) {
	return r.nodes, nil
}

func (r *queryResolver) GetNodes(ctx context.Context, first *int, after *string) (*model.Nodes, error) {
	// The cursor is base64 encoded by convention, so we need to decode it first
	var decodedCursor string
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, err
		}
		decodedCursor = string(b)
	}

	// Here we could query the DB to get data, e.g. SELECT * FROM messages WHERE chat_room_id = obj.ID AND timestamp < decodedCursor
	// Mocking for now
	nodes := make([]*model.Node1, *first)
	count := 0
	currentPage := false
	// If no cursor present start from the top
	if decodedCursor == "" {
		currentPage = true
	}
	hasNextPage := false

	// Iterating over the mocked messages to find the current page
	// In real world use-case you should fetch only the required part of data from the database
	for i, v := range r.nodes {
		node := v

		if currentPage && count < *first {
			nodes[count] = &model.Node1{
				Cursor: base64.StdEncoding.EncodeToString([]byte(v.IP)),
				Node:   node,
			}
			count++
		}

		if v.IP == decodedCursor {
			currentPage = true
		}

		// If there are any elements left after the current page we indicate that in the response
		if count == *first && i < len(r.nodes) {
			hasNextPage = true
		}
	}

	pageInfo := model.PageInfo{
		StartCursor: base64.StdEncoding.EncodeToString([]byte(nodes[0].Node.IP)),
		EndCursor:   base64.StdEncoding.EncodeToString([]byte(nodes[count-1].Node.IP)),
		HasNextPage: &hasNextPage,
	}

	fmt.Printf("MessagesConnection | first: %v, pageInfo: %+v \n", *first, pageInfo)

	mc := model.Nodes{
		Nodes:    nodes[:count],
		PageInfo: &pageInfo,
	}

	return &mc, nil
}

func (r *queryResolver) GetNodesByURL(ctx context.Context, sort *model.CovidDataSort) ([]*model.Covid, error) {
	r.covid = nil
	response, err := http.Get("https://data.covid19india.org/data.json")
	if err != nil {
		panic(err)
	}
	dataBytes, err1 := ioutil.ReadAll(response.Body)
	if err1 != nil {
		panic(err1)
	}
	// content := string(dataBytes)
	var myData Person
	json.Unmarshal(dataBytes, &myData)

	// json.Unmarshal([]byte(content), &myData)
	for _, val := range myData.Cases_time_series {
		cv := &model.Covid{
			Dailyconfirmed: val.Dailyconfirmed,
			Dailydeceased:  val.Dailydeceased,
			Dailyrecovered: val.Dailyrecovered,
			Date:           val.Date,
			Dateymd:        val.Dateymd,
			Totalconfirmed: val.Totalconfirmed,
			Totaldeceased:  val.Totaldeceased,
			Totalrecovered: val.Totalrecovered,
		}
		r.covid = append(r.covid, (*model.Covid)(cv))

	}
	sortRecords(r.covid, sort)
	defer response.Body.Close()
	return r.covid, nil
}

func (r *queryResolver) GetNodesByURLPagination(ctx context.Context, first *int, after *string) (*model.Covids, error) {
	response, err := http.Get("https://data.covid19india.org/data.json")
	if err != nil {
		panic(err)
	}
	dataBytes, err1 := ioutil.ReadAll(response.Body)
	if err1 != nil {
		panic(err1)
	}
	// content := string(dataBytes)
	var myData Person
	json.Unmarshal(dataBytes, &myData)
	// json.Unmarshal([]byte(content), &myData)
	for _, val := range myData.Cases_time_series {
		rand.Seed(time.Now().UnixNano())
		cv := &model.Covid{
			ID:             string(rune(rand.Int())),
			Dailyconfirmed: val.Dailyconfirmed,
			Dailydeceased:  val.Dailydeceased,
			Dailyrecovered: val.Dailyrecovered,
			Date:           val.Date,
			Dateymd:        val.Dateymd,
			Totalconfirmed: val.Totalconfirmed,
			Totaldeceased:  val.Totaldeceased,
			Totalrecovered: val.Totalrecovered,
		}
		r.covid = append(r.covid, (*model.Covid)(cv))

	}
	defer response.Body.Close()
	// return r.covid, nil
	// The cursor is base64 encoded by convention, so we need to decode it first
	var decodedCursor string
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, err
		}
		decodedCursor = string(b)
	}

	// Here we could query the DB to get data, e.g. SELECT * FROM messages WHERE chat_room_id = obj.ID AND timestamp < decodedCursor
	// Mocking for now
	covids := make([]*model.Covid10, *first)
	count := 0
	currentPage := false
	// If no cursor present start from the top
	if decodedCursor == "" {
		currentPage = true
	}
	hasNextPage := false

	// Iterating over the mocked messages to find the current page
	// In real world use-case you should fetch only the required part of data from the database
	for i, v := range r.covid {
		node := v

		if currentPage && count < *first {
			covids[count] = &model.Covid10{
				Cursor: base64.StdEncoding.EncodeToString([]byte(v.Dateymd)),
				Covid:  node,
			}
			count++
		}

		if v.Dateymd == decodedCursor {
			currentPage = true
		}

		// If there are any elements left after the current page we indicate that in the response
		if count == *first && i < len(r.covid) {
			hasNextPage = true
		}
	}

	pageInfo := model.PageInfo{
		StartCursor: base64.StdEncoding.EncodeToString([]byte(covids[0].Covid.Dateymd)),
		EndCursor:   base64.StdEncoding.EncodeToString([]byte(covids[count-1].Covid.Dateymd)),
		HasNextPage: &hasNextPage,
	}

	fmt.Printf("MessagesConnection | first: %v, pageInfo: %+v \n", *first, pageInfo)

	mc := model.Covids{
		Covids:   covids[:count],
		PageInfo: &pageInfo,
	}
	return &mc, nil
}

func sortRecords(records []*model.Covid, sortData *model.CovidDataSort) {
	// fmt.Println(records)
	// for _, el := range records {
	// 	fmt.Println(*el)
	// }
	// fmt.Println(sortData)
	var sorted map[string]interface{}
	marshal, err := json.Marshal(sortData)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	json.Unmarshal(([]byte)(marshal), &sorted)
	var sortkey string
	var sortvalue interface{}
	for key, val := range sorted {
		if val != nil {
			sortkey = key
			sortvalue = val
		}
	}
	// fmt.Println("Key values: ", sortkey, "\t", sortvalue)
	// if sortvalue == "asc" {
	sort.Slice(records, func(i, j int) bool {
		var sorted bool
		switch sortkey {
		case "totaldeceased":
			iRecord, ierr := strconv.Atoi(records[i].Totaldeceased)
			jRecord, jerr := strconv.Atoi(records[j].Totaldeceased)
			if sortvalue == "asc" {
				if ierr == nil && jerr == nil {
					sorted = iRecord < jRecord
				} else {
					sorted = records[i].Totaldeceased < records[j].Totaldeceased
				}
			} else if sortvalue == "desc" {
				if ierr == nil && jerr == nil {
					sorted = iRecord > jRecord
				} else {
					sorted = records[i].Totaldeceased > records[j].Totaldeceased
				}
			}
		case "totalconfirmed":
			iRecord, ierr := strconv.Atoi(records[i].Totalconfirmed)
			jRecord, jerr := strconv.Atoi(records[j].Totalconfirmed)
			if sortvalue == "asc" {
				if ierr == nil && jerr == nil {
					sorted = iRecord < jRecord
				} else {
					sorted = records[i].Totalconfirmed < records[j].Totalconfirmed
				}
			} else if sortvalue == "desc" {
				if ierr == nil && jerr == nil {
					sorted = iRecord > jRecord
				} else {
					sorted = records[i].Totalconfirmed > records[j].Totalconfirmed
				}
			}
		case "totalrecovered":
			iRecord, ierr := strconv.Atoi(records[i].Totalrecovered)
			jRecord, jerr := strconv.Atoi(records[j].Totalrecovered)
			if sortvalue == "asc" {
				if ierr == nil && jerr == nil {
					sorted = iRecord < jRecord
				} else {
					sorted = records[i].Totalrecovered < records[j].Totalrecovered
				}
			} else if sortvalue == "desc" {
				if ierr == nil && jerr == nil {
					sorted = iRecord > jRecord
				} else {
					sorted = records[i].Totalrecovered > records[j].Totalrecovered
				}
			}
		case "dailydeceased":
			iRecord, ierr := strconv.Atoi(records[i].Dailydeceased)
			jRecord, jerr := strconv.Atoi(records[j].Dailydeceased)
			if sortvalue == "asc" {
				if ierr == nil && jerr == nil {
					sorted = iRecord < jRecord
				} else {
					sorted = records[i].Dailydeceased < records[j].Dailydeceased
				}
			} else if sortvalue == "desc" {
				if ierr == nil && jerr == nil {
					sorted = iRecord > jRecord
				} else {
					sorted = records[i].Dailydeceased > records[j].Dailydeceased
				}
			}
		case "dailyconfirmed":
			iRecord, ierr := strconv.Atoi(records[i].Dailyconfirmed)
			jRecord, jerr := strconv.Atoi(records[j].Dailyconfirmed)
			if sortvalue == "asc" {
				if ierr == nil && jerr == nil {
					sorted = iRecord < jRecord
				} else {
					sorted = records[i].Dailyconfirmed < records[j].Dailyconfirmed
				}
			} else if sortvalue == "desc" {
				if ierr == nil && jerr == nil {
					sorted = iRecord > jRecord
				} else {
					sorted = records[i].Dailyconfirmed > records[j].Dailyconfirmed
				}
			}
		case "dailyrecovered":
			iRecord, ierr := strconv.Atoi(records[i].Dailyrecovered)
			jRecord, jerr := strconv.Atoi(records[j].Dailyrecovered)
			if sortvalue == "asc" {
				if ierr == nil && jerr == nil {
					sorted = iRecord < jRecord
				} else {
					sorted = records[i].Dailyrecovered < records[j].Dailyrecovered
				}
			} else if sortvalue == "desc" {
				if ierr == nil && jerr == nil {
					sorted = iRecord > jRecord
				} else {
					sorted = records[i].Dailyrecovered > records[j].Dailyrecovered
				}
			}
		case "date":
			if sortvalue == "asc" {
				sorted = records[i].Dateymd < records[j].Dateymd
			} else if sortvalue == "desc" {
				sorted = records[i].Dateymd > records[j].Dateymd
			}
		}
		return sorted
	})
	// } else if sortvalue == "desc" {
	// 	sort.Slice()
	// }
}

func (r *subscriptionResolver) VideoAdded(ctx context.Context, repoFullName string) (<-chan *model.Node, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
type Covid struct {
	Dailyconfirmed string `json:"dailyconfirmed"`
	Dailydeceased  string `json:"dailydeceased"`
	Dailyrecovered string `json:"dailyrecovered"`
	Date           string `json:"date"`
	Dateymd        string `json:"dateymd"`
	Totalconfirmed string `json:"totalconfirmed"`
	Totaldeceased  string `json:"totaldeceased"`
	Totalrecovered string `json:"totalrecovered"`
}
type Covid_1 struct {
	Active          string `json:"active"`
	Confirmed       string `json:"confirmed"`
	Deaths          string `json:"deaths"`
	Deltaconfirmed  string `json:"deltaconfirmed"`
	Deltadeaths     string `json:"deltadeaths"`
	Deltarecovered  string `json:"deltarecovered"`
	Lastupdatedtime string `json:"lastupdatedtime"`
	Migratedother   string `json:"migratedother"`
	Recovered       string `json:"recovered"`
	State           string `json:"state"`
	Statecode       string `json:"statecode"`
	Statenotes      string `json:"statenotes"`
}
type Covid_2 struct {
	Firstdoseadministered string `json:"firstdoseadministered"`
}
type Person struct {
	Cases_time_series []Covid   `json:"cases_time_series"`
	Statewise         []Covid_1 `json:"statewise"`
	Tested            []Covid_2 `json:"tested"`
}
