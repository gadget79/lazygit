package context

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

type TagsContext struct {
	*TagsViewModel
	*BaseContext
	*ListContextTrait
}

var _ types.IListContext = (*TagsContext)(nil)

func NewTagsContext(
	getModel func() []*models.Tag,
	getView func() *gocui.View,
	getDisplayStrings func(startIdx int, length int) [][]string,

	onFocus func(...types.OnFocusOpts) error,
	onRenderToMain func(...types.OnFocusOpts) error,
	onFocusLost func() error,

	c *types.ControllerCommon,
) *TagsContext {
	baseContext := NewBaseContext(NewBaseContextOpts{
		ViewName:   "branches",
		WindowName: "branches",
		Key:        TAGS_CONTEXT_KEY,
		Kind:       types.SIDE_CONTEXT,
	})

	self := &TagsContext{}
	takeFocus := func() error { return c.PushContext(self) }

	list := NewTagsViewModel(getModel)
	viewTrait := NewViewTrait(getView)
	listContextTrait := &ListContextTrait{
		base:      baseContext,
		listTrait: list.ListTrait,
		viewTrait: viewTrait,

		GetDisplayStrings: getDisplayStrings,
		OnFocus:           onFocus,
		OnRenderToMain:    onRenderToMain,
		OnFocusLost:       onFocusLost,
		takeFocus:         takeFocus,

		// TODO: handle this in a trait
		RenderSelection: false,

		c: c,
	}

	self.BaseContext = baseContext
	self.ListContextTrait = listContextTrait
	self.TagsViewModel = list

	return self
}

type TagsViewModel struct {
	*ListTrait
	getModel func() []*models.Tag
}

func (self *TagsViewModel) GetItemsLength() int {
	return len(self.getModel())
}

func (self *TagsViewModel) GetSelectedTag() *models.Tag {
	if self.GetItemsLength() == 0 {
		return nil
	}

	return self.getModel()[self.GetSelectedLineIdx()]
}

func (self *TagsViewModel) GetSelectedItem() (types.ListItem, bool) {
	item := self.GetSelectedTag()
	return item, item != nil
}

func NewTagsViewModel(getModel func() []*models.Tag) *TagsViewModel {
	self := &TagsViewModel{
		getModel: getModel,
	}

	self.ListTrait = &ListTrait{
		selectedIdx: 0,
		HasLength:   self,
	}

	return self
}

func clamp(x int, min int, max int) int {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}