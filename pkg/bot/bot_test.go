package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestCoinFlip(t *testing.T) {
	fmt.Println("Flipping...")

	time.Sleep(3 * time.Second)
	fmt.Println("...")

	for i := 0; i < 5; i++ {
		time.Sleep(3 * time.Second)
		x1 := rand.NewSource(time.Now().UnixNano())
		y1 := rand.New(x1)
		randNum := y1.Intn(200)

		if randNum%2 == 0 {
			fmt.Println("It landed heads")

		} else {
			fmt.Println("It landed tails")
		}
	}
}

func TestMemberHasRole(t *testing.T) {
	roleName := "test"
	s := discordgo.NewState()

	member, err := s.Member("293416960237240320", "289217573004902400")
	if err != nil {
		t.Fatal(err)
	}

	// memberRoles := make([]string, len(member.Roles))

	for _, role := range member.Roles {
		if role == "@everyone" {
			continue
		}

		if strings.ToLower(role) == roleName {
			fmt.Println("Role not found")
		}
	}
}
