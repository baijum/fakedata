package main

import (
	"flag"
	"fmt"
	"math/rand"

	"github.com/dmgk/faker"
	"github.com/kaaryasthan/kaaryasthan/config"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/item/model"
	"github.com/kaaryasthan/kaaryasthan/label/model"
	"github.com/kaaryasthan/kaaryasthan/milestone/model"
	"github.com/kaaryasthan/kaaryasthan/project/model"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// Profile represents the dataset based volume of data
type Profile struct {
	Users       int
	Projects    int
	Milestones  int
	Labels      int
	Items       int
	Discussions int
	Comments    int
}

var profileList = map[string]Profile{
	"nano": {
		Users:       1,
		Projects:    1,
		Milestones:  1,
		Labels:      1,
		Items:       2,
		Discussions: 3,
		Comments:    4,
	},
	"micro": {
		Users:       2,
		Projects:    2,
		Milestones:  4,
		Labels:      5,
		Items:       10,
		Discussions: 30,
		Comments:    70,
	},
	"tiny": {
		Users:       4,
		Projects:    3,
		Milestones:  8,
		Labels:      7,
		Items:       100,
		Discussions: 300,
		Comments:    700,
	},
	"small": {
		Users:       10,
		Projects:    6,
		Milestones:  20,
		Labels:      25,
		Items:       1000,
		Discussions: 3000,
		Comments:    7000,
	},
	"medium": {
		Users:       50,
		Projects:    12,
		Milestones:  100,
		Labels:      125,
		Items:       10000,
		Discussions: 30000,
		Comments:    70000,
	},
	"large": {
		Users:       250,
		Projects:    24,
		Milestones:  500,
		Labels:      625,
		Items:       100000,
		Discussions: 300000,
		Comments:    700000,
	},
	"huge": {
		Users:       1250,
		Projects:    48,
		Milestones:  2500,
		Labels:      3125,
		Items:       1000000,
		Discussions: 3000000,
		Comments:    7000000,
	},
}

var profile = flag.String("profile", "", `Name of dataset profile.
	Available profiles: nano, micro, tiny, small, medium, large, hugea`)

func generateData(p Profile) {
	conf := config.Config.PostgresConfig()
	DB := db.Connect(conf)
	userList := make([]user.User, 0, p.Users)
	for i := 0; i < p.Users; i++ {
		usrDS := user.NewDatastore(DB)
		usr := &user.User{Username: faker.Internet().UserName(), Name: faker.Name().Name(), Email: faker.Internet().SafeEmail(), Password: "Secret@123"}
		if err := usrDS.Create(usr); err != nil {
			fmt.Println(err)
		}
		userList = append(userList, *usr)
	}

	projectList := make([]project.Project, 0, p.Projects)
	for i := 0; i < p.Projects; i++ {
		prjDS := project.NewDatastore(DB)
		prj := &project.Project{Name: faker.Internet().UserName(), Description: faker.Commerce().ProductName()}
		nu := p.Users
		if nu > 3 {
			nu = 3
		}
		usr := userList[rand.Intn(nu)]
		if err := prjDS.Create(&usr, prj); err != nil {
			fmt.Println(err)
		}
		projectList = append(projectList, *prj)
	}

	for i := 0; i < p.Milestones; i++ {
		prj := projectList[rand.Intn(p.Projects)]
		milDS := milestone.NewDatastore(DB)
		mil := &milestone.Milestone{Name: faker.Internet().UserName(), Description: faker.Commerce().ProductName(), ProjectID: prj.ID}
		usr := userList[rand.Intn(p.Users)]
		if err := milDS.Create(&usr, mil); err != nil {
			fmt.Println(err)
		}
	}

	for i := 0; i < p.Labels; i++ {
		prj := projectList[rand.Intn(p.Projects)]
		lbl := &label.Label{Name: faker.Internet().UserName(), Color: "#ee0701", ProjectID: prj.ID}
		usr := userList[rand.Intn(p.Users)]
		lblDS := label.NewDatastore(DB)
		if err := lblDS.Create(&usr, lbl); err != nil {
			fmt.Println(err)
		}
	}

	itemList := make([]item.Item, 0, p.Items)
	for i := 0; i < p.Items; i++ {
		prj := projectList[rand.Intn(p.Projects)]
		itm := &item.Item{Title: faker.Internet().UserName(), Description: faker.Hacker().SaySomethingSmart(), ProjectID: prj.ID}
		usr := userList[rand.Intn(p.Users)]
		itmDS := item.NewDatastore(DB)
		if err := itmDS.Create(&usr, itm); err != nil {
			fmt.Println(err)
		}
		itemList = append(itemList, *itm)
	}

	discussionList := make([]item.Discussion, 0, p.Discussions)
	for i := 0; i < p.Discussions; i++ {
		itm := itemList[rand.Intn(p.Items)]
		disc := &item.Discussion{Body: faker.Hacker().SaySomethingSmart(), ItemID: itm.ID}
		usr := userList[rand.Intn(p.Users)]
		discDS := item.NewDiscussionDatastore(DB)
		if err := discDS.Create(&usr, disc); err != nil {
			fmt.Println(err)
		}
		discussionList = append(discussionList, *disc)
	}

	for i := 0; i < p.Comments; i++ {
		disc := discussionList[rand.Intn(p.Discussions)]
		com := &item.Comment{Body: faker.Hacker().SaySomethingSmart(), DiscussionID: disc.ID}
		usr := userList[rand.Intn(p.Users)]
		comDS := item.NewCommentDatastore(DB)
		if err := comDS.Create(&usr, com); err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	flag.Parse()

	if *profile != "" {
		p, ok := profileList[*profile]
		if !ok {
			fmt.Println("Profile not found:", *profile)
			return
		}
		generateData(p)
	}
}
