package tap

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/pkg/utils"
	"github.com/zeusWPI/scc/tui/view"
)

func updateOrders(ctx context.Context, view view.View) (tea.Msg, error) {
	m := view.(*Model)

	lastOrder, err := m.repo.GetLast(ctx)
	if err != nil {
		return nil, err
	}
	if lastOrder == nil {
		lastOrder = &model.Tap{OrderID: -1}
	}

	if lastOrder.OrderID == m.lastOrderID {
		return nil, nil
	}

	counts, err := m.repo.GetCountByCategory(ctx)
	if err != nil {
		return nil, nil
	}

	return Msg{items: utils.SliceDereference(counts)}, nil
}
