package etic

type Cartridge interface {
	Title() string
	Size() Point[float32]
	BOOT(*Console)
	TIC(*Console)
}
type Cart struct {
	title  string
	size   Point[float32]
	OnBoot func(*Console)
	OnTic  func(*Console)
}

func NewCart(title string, w, h float32) *Cart {
	return &Cart{title: title,
		size: Pt(w, h),
	}
}
func (c *Cart) Title() string        { return c.title }
func (c *Cart) Size() Point[float32] { return c.size }
func (c *Cart) BOOT(t *Console) {
	if c.OnBoot != nil {
		c.OnBoot(t)
	}
}
func (c *Cart) TIC(t *Console) {
	if c.OnTic != nil {
		c.OnTic(t)
	}
}
