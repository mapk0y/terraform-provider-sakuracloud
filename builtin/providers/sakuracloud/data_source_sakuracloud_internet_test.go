package sakuracloud

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sacloud/libsacloud/api"
	"testing"
)

func TestAccSakuraCloudInternetDataSource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudInternetDataSourceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceInternetBase,
				Check:  testAccCheckSakuraCloudInternetDataSourceID("sakuracloud_internet.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceInternetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetDataSourceID("data.sakuracloud_internet.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "name", "name_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "tags.2", "tag3"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "nw_mask_len", "28"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "band_width", "100"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "server_ids.#", "0"),
					resource.TestCheckResourceAttr("data.sakuracloud_internet.foobar", "nw_ipaddresses.#", "11"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceInternetConfig_With_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetDataSourceID("data.sakuracloud_internet.foobar"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceInternetConfig_NotExists,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetDataSourceNotExists("data.sakuracloud_internet.foobar"),
				),
			},
			{
				Destroy: true,
				Config:  testAccCheckSakuraCloudDataSourceInternetConfig_With_NotExists_Tag,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudInternetDataSourceNotExists("data.sakuracloud_internet.foobar"),
				),
			},
		},
	})
}

func testAccCheckSakuraCloudInternetDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Internet data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Internet data source ID not set")
		}
		return nil
	}
}

func testAccCheckSakuraCloudInternetDataSourceNotExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if ok {
			return fmt.Errorf("Found Internet data source: %s", n)
		}
		return nil
	}
}

func testAccCheckSakuraCloudInternetDataSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sakuracloud_internet" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		_, err := client.Internet.Read(toSakuraCloudID(rs.Primary.ID))

		if err == nil {
			return errors.New("Internet still exists")
		}
	}

	return nil
}

var testAccCheckSakuraCloudDataSourceInternetBase = `
resource "sakuracloud_internet" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
`

var testAccCheckSakuraCloudDataSourceInternetConfig = `
resource "sakuracloud_internet" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_internet" "foobar" {
    filter = {
	name = "Name"
	values = ["name_test"]
    }
}`

var testAccCheckSakuraCloudDataSourceInternetConfig_With_Tag = `
resource "sakuracloud_internet" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_internet" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1","tag3"]
    }
}`

var testAccCheckSakuraCloudDataSourceInternetConfig_With_NotExists_Tag = `
resource "sakuracloud_internet" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_internet" "foobar" {
    filter = {
	name = "Tags"
	values = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
    }
}`

var testAccCheckSakuraCloudDataSourceInternetConfig_NotExists = `
resource "sakuracloud_internet" "foobar" {
    name = "name_test"
    description = "description_test"
    tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_internet" "foobar" {
    filter = {
	name = "Name"
	values = ["xxxxxxxxxxxxxxxxxx"]
    }
}`
