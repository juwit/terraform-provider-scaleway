package scaleway

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccScalewayDataSourceVPCGatewayNetwork_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayVPCGatewayNetworkDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "scaleway_vpc_private_network" "pn01" {
					}
					
					resource "scaleway_vpc_public_gateway_ip" "gw01" {
					}
					
					resource "scaleway_vpc_public_gateway_dhcp" "dhcp01" {
					  subnet = "192.168.1.0/24"
					  push_default_route = true
					}
					
					resource "scaleway_vpc_public_gateway" "pg01" {
					  type = "VPC-GW-S"
					  ip_id = scaleway_vpc_public_gateway_ip.gw01.id
					}
					
					resource "scaleway_vpc_gateway_network" "main" {
					  gateway_id = scaleway_vpc_public_gateway.pg01.id
					  private_network_id = scaleway_vpc_private_network.pn01.id
					  dhcp_id = scaleway_vpc_public_gateway_dhcp.dhcp01.id
					  cleanup_dhcp       = true
					  enable_masquerade  = true
					}`,
			},
			{
				Config: `
					resource "scaleway_vpc_private_network" "pn01" {
					}
					
					resource "scaleway_vpc_public_gateway_ip" "gw01" {
					}
					
					resource "scaleway_vpc_public_gateway_dhcp" "dhcp01" {
					  subnet = "192.168.1.0/24"
					  push_default_route = true
					}
					
					resource "scaleway_vpc_public_gateway" "pg01" {
					  type = "VPC-GW-S"
					  ip_id = scaleway_vpc_public_gateway_ip.gw01.id
					}
					
					resource "scaleway_vpc_gateway_network" "main" {
					  gateway_id = scaleway_vpc_public_gateway.pg01.id
					  private_network_id = scaleway_vpc_private_network.pn01.id
					  dhcp_id = scaleway_vpc_public_gateway_dhcp.dhcp01.id
					  cleanup_dhcp       = true
					  enable_masquerade  = true
					}

					data scaleway_vpc_gateway_network by_id {
						gateway_network_id = scaleway_vpc_gateway_network.main.id
					}

					data scaleway_vpc_gateway_network by_gateway_and_pn {
						gateway_id = scaleway_vpc_public_gateway.pg01.id
						private_network_id = scaleway_vpc_private_network.pn01.id
					}
					`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayVPCGatewayNetworkExists(tt, "scaleway_vpc_gateway_network.main"),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_gateway_network.by_id", "id",
						"scaleway_vpc_gateway_network.main", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.scaleway_vpc_gateway_network.by_gateway_and_pn", "id",
						"scaleway_vpc_gateway_network.main", "id",
					),
				),
			},
		},
	})
}
