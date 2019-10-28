package lsp

import (
	"context"

	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/source"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/span"
	errors "golang.org/x/xerrors"
)

func (s *Server) executeCommand(ctx context.Context, params *protocol.ExecuteCommandParams) (interface{}, error) {
	switch params.Command {
	case "tidy":
		if len(params.Arguments) == 0 || len(params.Arguments) > 1 {
			return nil, errors.Errorf("expected one file URI for call to `go mod tidy`, got %v", params.Arguments)
		}
		// Confirm that this action is being taken on a go.mod file.
		uri := span.NewURI(params.Arguments[0].(string))
		view := s.session.ViewOf(uri)
		f, err := view.GetFile(ctx, uri)
		if err != nil {
			return nil, err
		}
		fh := view.Snapshot().Handle(ctx, f)
		if fh.Identity().Kind != source.Mod {
			return nil, errors.Errorf("%s is not a mod file", uri)
		}
		// Run go.mod tidy on the view.
		if err := source.ModTidy(ctx, view); err != nil {
			return nil, err
		}
	}
	return nil, nil
}