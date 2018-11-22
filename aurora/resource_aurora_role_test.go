package aurora

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"testing"

	r "github.com/hashicorp/terraform/helper/resource"
)

func TestAuroraRole_basic(t *testing.T) {
	role := acctest.RandomWithPrefix("test")

	r.Test(t, r.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: testAccAuroraRole_basic(role),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("aurora_role.test", "name", role),
					r.TestCheckResourceAttr("aurora_role.test", "cpu", "20.5"),
					r.TestCheckResourceAttr("aurora_role.test", "ram_mb", "512"),
					r.TestCheckResourceAttr("aurora_role.test", "disk_mb", "1024"),
				),
			},
		},
	})
}

func testAccAuroraRole_basic(role string) string {
	return fmt.Sprintf(`
resource "aurora_role" "test" {
  name = "%s"
  cpu = 20.5
  ram_mb = 512
  disk_mb = 1024
}

`, role)
}
