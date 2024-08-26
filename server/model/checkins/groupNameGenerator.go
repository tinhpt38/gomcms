package checkins

import "fmt"

type GroupNameGenerator struct {
	NumericCount    int
	AlphabeticCount int
	ExistingNames   map[string]bool
}

func (g *GroupNameGenerator) GenerateName(nameType string) string {
	switch nameType {
	case "baseOnNumberic":
		{
			for {
				name := fmt.Sprintf("Nhóm %d", g.NumericCount)
				if !g.ExistingNames[name] {
					g.ExistingNames[name] = true
					g.NumericCount++
					return name
				}
				g.NumericCount++
			}

			break
		}
	case "baseOnAlphabet":
		{
			for {
				name := fmt.Sprintf("Nhóm %c", 'A'+g.AlphabeticCount)
				if !g.ExistingNames[name] {
					g.ExistingNames[name] = true
					g.AlphabeticCount++
					return name
				}
				g.AlphabeticCount++
			}

			break
		}
	}
	return ""
}
