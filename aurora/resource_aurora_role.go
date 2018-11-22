package aurora

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceAuroraRole() *schema.Resource {
	return &schema.Resource{

		Create: quotaCreate,
		Read:   quotaRead,
		Update: quotaCreate,
		Delete: quotaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cpu": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"ram_mb": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"disk_mb": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func quotaCreate(d *schema.ResourceData, m interface{}) error {
	var name string
	var c float64
	var r int64
	var dsk int64

	config := m.(*Config)

	name = d.Get("name").(string)
	cpu := d.Get("cpu").(float64)
	c = float64(cpu)
	ram := d.Get("ram_mb").(int)
	r = int64(ram)
	disk := d.Get("disk_mb").(int)
	dsk = int64(disk)

	log.Printf("[DEBUG] creating quota for role: %s", name)

	if _, err := config.client.SetQuota(name, &c, &r, &dsk); err != nil {
		fmt.Print(err)
		log.Printf("[DEBUG] error with setting quota: %s", err)
		return err
	}
	if d.Id() != name {
		d.SetId(name)
	}
	return quotaRead(d, m)
}

func quotaRead(d *schema.ResourceData, m interface{}) error {

	name := d.Id()
	config := m.(*Config)
	if res, err := config.client.GetQuota(name); err != nil {
		d.SetId("")
		log.Printf("[DEBUG] error with reading quota: %s", err)
		return err
	} else {
		for key := range res.Result_.GetQuotaResult_.GetQuota().GetResources() {
			switch {
			case key.IsSetNumCpus():
				cpu := key.GetNumCpus()
				if d.Get("cpu").(float64) != float64(cpu) {
					d.Set("cpu", float64(cpu))
				}
			case key.IsSetRamMb():
				ram := key.GetRamMb()
				if d.Get("ram_mb").(int) != int(ram) {
					d.Set("ram_mb", int(ram))
				}
			case key.IsSetDiskMb():
				disk := key.GetDiskMb()
				if d.Get("disk_mb").(int) != int(disk) {
					d.Set("disk_mb", int(disk))
				}
			}
		}
	}
	return nil
}

func quotaDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	var cpu float64 = 0
	var ram int64 = 0
	var disk int64 = 0
	name := d.Id()
	if _, err := config.client.SetQuota(name, &cpu, &ram, &disk); err != nil {
		log.Printf("[DEBUG] error with deleting quota: %s", err)
		return err
	}
	return nil
}
