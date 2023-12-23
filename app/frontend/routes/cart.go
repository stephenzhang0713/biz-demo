package routes

import (
	"context"
	"fmt"
	"github.com/baiyutang/gomall/app/frontend/kitex_gen/product"
	"strconv"

	"github.com/baiyutang/gomall/app/frontend/infra/rpc"
	"github.com/baiyutang/gomall/app/frontend/kitex_gen/cart"
	frontendutils "github.com/baiyutang/gomall/app/frontend/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func RegisterCart(h *server.Hertz) {
	h.GET("/cart", func(ctx context.Context, c *app.RequestContext) {
		var items []map[string]string
		carts, err := rpc.CartClient.GetCart(ctx, &cart.GetCartRequest{UserId: uint32(ctx.Value(frontendutils.UserIdKey).(int))})
		if err != nil {
			c.JSON(500, "get cart error")
		}
		for _, v := range carts.Items {
			p, err := rpc.ProductClient.GetProduct(ctx, &product.GetProductRequest{Id: v.GetProductId()})
			if err != nil {
				continue
			}
			items = append(items, map[string]string{"Name": p.Name, "Description": p.Description, "Picture": p.Picture, "Price": strconv.FormatFloat(float64(p.Price), 'f', 2, 64), "Qty": strconv.Itoa(int(v.Quantity))})
		}
		c.HTML(consts.StatusOK, "cart", frontendutils.WarpResponse(ctx, c, utils.H{
			"title":    "Cart",
			"items":    items,
			"cart_num": 10,
		}))
	})
	type form struct {
		ProductId  uint32 `json:"productId" form:"productId"`
		ProductNum uint32 `json:"productNum" form:"productNum"`
	}
	h.POST("/cart", func(ctx context.Context, c *app.RequestContext) {
		var f form
		c.BindAndValidate(&f)

		r, err := rpc.CartClient.AddItem(ctx, &cart.AddItemRequest{UserId: uint32(ctx.Value(frontendutils.UserIdKey).(int)), Item: &cart.CartItem{
			ProductId: f.ProductId,
			Quantity:  int32(f.ProductNum),
		}})
		fmt.Println(r, err)
		c.Redirect(consts.StatusFound, []byte("/cart"))
	})
}
