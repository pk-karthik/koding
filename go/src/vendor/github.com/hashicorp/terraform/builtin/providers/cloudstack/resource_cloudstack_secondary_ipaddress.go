package cloudstack

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/xanzy/go-cloudstack/cloudstack"
)

func resourceCloudStackSecondaryIPAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackSecondaryIPAddressCreate,
		Read:   resourceCloudStackSecondaryIPAddressRead,
		Delete: resourceCloudStackSecondaryIPAddressDelete,

		Schema: map[string]*schema.Schema{
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"ipaddress": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Please use the `ip_address` field instead",
			},

			"nic_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"nicid": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Please use the `nic_id` field instead",
			},

			"virtual_machine_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"virtual_machine": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Please use the `virtual_machine_id` field instead",
			},
		},
	}
}

func resourceCloudStackSecondaryIPAddressCreate(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	nicid, ok := d.GetOk("nic_id")
	if !ok {
		nicid, ok = d.GetOk("nicid")
	}
	if !ok {
		virtualmachine, ok := d.GetOk("virtual_machine_id")
		if !ok {
			virtualmachine, ok = d.GetOk("virtual_machine")
		}
		if !ok {
			return errors.New(
				"Either `virtual_machine_id` or [deprecated] `virtual_machine` must be provided.")
		}

		// Retrieve the virtual_machine ID
		virtualmachineid, e := retrieveID(cs, "virtual_machine", virtualmachine.(string))
		if e != nil {
			return e.Error()
		}

		// Get the virtual machine details
		vm, count, err := cs.VirtualMachine.GetVirtualMachineByID(virtualmachineid)
		if err != nil {
			if count == 0 {
				log.Printf("[DEBUG] Virtual Machine %s does no longer exist", virtualmachineid)
				d.SetId("")
				return nil
			}
			return err
		}

		nicid = vm.Nic[0].Id
	}

	// Create a new parameter struct
	p := cs.Nic.NewAddIpToNicParams(nicid.(string))

	// If there is a ipaddres supplied, add it to the parameter struct
	ipaddress, ok := d.GetOk("ip_address")
	if !ok {
		ipaddress, ok = d.GetOk("ipaddress")
	}
	if ok {
		p.SetIpaddress(ipaddress.(string))
	}

	ip, err := cs.Nic.AddIpToNic(p)
	if err != nil {
		return err
	}

	d.SetId(ip.Id)

	return nil
}

func resourceCloudStackSecondaryIPAddressRead(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	virtualmachine, ok := d.GetOk("virtual_machine_id")
	if !ok {
		virtualmachine, ok = d.GetOk("virtual_machine")
	}
	if !ok {
		return errors.New(
			"Either `virtual_machine_id` or [deprecated] `virtual_machine` must be provided.")
	}

	// Retrieve the virtual_machine ID
	virtualmachineid, e := retrieveID(cs, "virtual_machine", virtualmachine.(string))
	if e != nil {
		return e.Error()
	}

	// Get the virtual machine details
	vm, count, err := cs.VirtualMachine.GetVirtualMachineByID(virtualmachineid)
	if err != nil {
		if count == 0 {
			log.Printf("[DEBUG] Virtual Machine %s does no longer exist", virtualmachineid)
			d.SetId("")
			return nil
		}
		return err
	}

	nicid, ok := d.GetOk("nic_id")
	if !ok {
		nicid, ok = d.GetOk("nicid")
	}
	if !ok {
		nicid = vm.Nic[0].Id
	}

	p := cs.Nic.NewListNicsParams(virtualmachineid)
	p.SetNicid(nicid.(string))

	l, err := cs.Nic.ListNics(p)
	if err != nil {
		return err
	}

	if l.Count == 0 {
		log.Printf("[DEBUG] NIC %s does no longer exist", d.Get("nicid").(string))
		d.SetId("")
		return nil
	}

	if l.Count > 1 {
		return fmt.Errorf("Found more then one possible result: %v", l.Nics)
	}

	for _, ip := range l.Nics[0].Secondaryip {
		if ip.Id == d.Id() {
			d.Set("ip_address", ip.Ipaddress)
			d.Set("nic_id", l.Nics[0].Id)
			d.Set("virtual_machine_id", l.Nics[0].Virtualmachineid)
			return nil
		}
	}

	log.Printf("[DEBUG] IP %s no longer exist", d.Get("ip_address").(string))
	d.SetId("")

	return nil
}

func resourceCloudStackSecondaryIPAddressDelete(d *schema.ResourceData, meta interface{}) error {
	cs := meta.(*cloudstack.CloudStackClient)

	// Create a new parameter struct
	p := cs.Nic.NewRemoveIpFromNicParams(d.Id())

	log.Printf("[INFO] Removing secondary IP address: %s", d.Get("ip_address").(string))
	if _, err := cs.Nic.RemoveIpFromNic(p); err != nil {
		// This is a very poor way to be told the ID does no longer exist :(
		if strings.Contains(err.Error(), fmt.Sprintf(
			"Invalid parameter id value=%s due to incorrect long value format, "+
				"or entity does not exist", d.Id())) {
			return nil
		}

		return fmt.Errorf("Error removing secondary IP address: %s", err)
	}

	return nil
}
