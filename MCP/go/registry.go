package main

import (
	"github.com/solar-vps/mcp-server/config"
	"github.com/solar-vps/mcp-server/models"
	tools_tickets "github.com/solar-vps/mcp-server/tools/tickets"
	tools_pods "github.com/solar-vps/mcp-server/tools/pods"
	tools_dns "github.com/solar-vps/mcp-server/tools/dns"
	tools_domains "github.com/solar-vps/mcp-server/tools/domains"
	tools_solarray "github.com/solar-vps/mcp-server/tools/solarray"
	tools_contacts "github.com/solar-vps/mcp-server/tools/contacts"
	tools_key "github.com/solar-vps/mcp-server/tools/key"
)

func GetAll(cfg *config.APIConfig) []models.Tool {
	return []models.Tool{
		tools_tickets.CreateGet_ticketsTool(cfg),
		tools_tickets.CreatePost_tickets_department_addTool(cfg),
		tools_pods.CreateGet_pods_podidTool(cfg),
		tools_pods.CreateGet_pods_podid_pingTool(cfg),
		tools_tickets.CreateGet_tickets_ticketidTool(cfg),
		tools_tickets.CreatePost_tickets_ticketid_updateTool(cfg),
		tools_dns.CreateGet_dns_domainTool(cfg),
		tools_domains.CreatePost_domains_addTool(cfg),
		tools_solarray.CreateGet_solarray_criticalTool(cfg),
		tools_contacts.CreateGet_contactsTool(cfg),
		tools_dns.CreatePost_dns_domain_addTool(cfg),
		tools_key.CreateGet_key_generateTool(cfg),
		tools_key.CreateGet_key_getTool(cfg),
		tools_pods.CreateGet_pods_podid_actionTool(cfg),
		tools_domains.CreatePost_domains_deleteTool(cfg),
		tools_dns.CreatePost_dns_domain_deleteTool(cfg),
		tools_domains.CreateGet_domainsTool(cfg),
		tools_dns.CreatePost_dns_domain_updateTool(cfg),
		tools_pods.CreateGet_podsTool(cfg),
		tools_solarray.CreateGet_solarrayTool(cfg),
	}
}
