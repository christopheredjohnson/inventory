package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Item struct {
	Name      string
	IconRect  rl.Rectangle // from spritesheet
	Stackable bool
}

type ItemSlot struct {
	Item  *Item
	Count int
}

type Inventory struct {
	Slots        []ItemSlot
	Cols         int
	Rows         int
	SlotSize     int
	SlotPadding  int
	Position     rl.Vector2
	ItemTexture  rl.Texture2D
	HeldItem     *ItemSlot // for dragging
	HeldFromSlot int       // index of source slot
}

func (inv *Inventory) AddItem(item *Item, count int) {
	for i := range inv.Slots {
		slot := &inv.Slots[i]
		if slot.Item != nil && slot.Item.Name == item.Name && item.Stackable {
			slot.Count += count
			return
		}
	}
	for i := range inv.Slots {
		if inv.Slots[i].Item == nil {
			inv.Slots[i] = ItemSlot{Item: item, Count: count}
			return
		}
	}
}
func (inv *Inventory) Draw() {
	for i, slot := range inv.Slots {
		col := i % inv.Cols
		row := i / inv.Cols
		x := inv.Position.X + float32(col*(inv.SlotSize+inv.SlotPadding))
		y := inv.Position.Y + float32(row*(inv.SlotSize+inv.SlotPadding))
		rect := rl.NewRectangle(x, y, float32(inv.SlotSize), float32(inv.SlotSize))

		rl.DrawRectangleRec(rect, rl.Fade(rl.LightGray, 0.8))
		rl.DrawRectangleLinesEx(rect, 1, rl.Black)

		if slot.Item != nil {
			iconW := slot.Item.IconRect.Width
			iconH := slot.Item.IconRect.Height
			offsetX := (float32(inv.SlotSize) - iconW) / 2
			offsetY := (float32(inv.SlotSize) - iconH) / 2
			destPos := rl.NewVector2(x+offsetX, y+offsetY)

			rl.DrawTextureRec(inv.ItemTexture, slot.Item.IconRect, destPos, rl.White)

			if slot.Item.Stackable && slot.Count > 1 {
				rl.DrawText(fmt.Sprintf("%d", slot.Count), int32(x+2), int32(y+2), 12, rl.Yellow)
			}
		}
	}

	if inv.HeldItem != nil {
		pos := rl.GetMousePosition()
		rl.DrawTextureRec(inv.ItemTexture, inv.HeldItem.Item.IconRect, pos, rl.White)
	}
}

func (inv *Inventory) Update() {
	mouse := rl.GetMousePosition()
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		for i, _ := range inv.Slots {
			col := i % inv.Cols
			row := i / inv.Cols
			x := inv.Position.X + float32(col*(inv.SlotSize+inv.SlotPadding))
			y := inv.Position.Y + float32(row*(inv.SlotSize+inv.SlotPadding))
			rect := rl.NewRectangle(x, y, float32(inv.SlotSize), float32(inv.SlotSize))

			if rl.CheckCollisionPointRec(mouse, rect) {
				inv.handleSlotClick(i)
				break
			}
		}
	}
}

func (inv *Inventory) handleSlotClick(i int) {
	slot := &inv.Slots[i]

	// Pick up item
	if inv.HeldItem == nil && slot.Item != nil {
		inv.HeldItem = &ItemSlot{Item: slot.Item, Count: slot.Count}
		inv.HeldFromSlot = i
		inv.Slots[i] = ItemSlot{}
		return
	}

	// Drop item
	if inv.HeldItem != nil {
		target := &inv.Slots[i]
		if target.Item == nil {
			inv.Slots[i] = *inv.HeldItem
			inv.HeldItem = nil
		} else if target.Item.Name == inv.HeldItem.Item.Name && target.Item.Stackable {
			target.Count += inv.HeldItem.Count
			inv.HeldItem = nil
		} else {
			// Swap
			inv.Slots[inv.HeldFromSlot] = *target
			inv.Slots[i] = *inv.HeldItem
			inv.HeldItem = nil
		}
	}
}
