package filters

import (
	"reflect"

	"github.com/juju/errors"

	"github.com/moiot/gravity/pkg/config"
	"github.com/moiot/gravity/pkg/core"
	"github.com/moiot/gravity/pkg/registry"
)

func NewFilters(filterConfigs []config.GenericPluginConfig) ([]core.IFilter, error) {
	var retFilters []core.IFilter
	for _, c := range filterConfigs {
		if c.Type == "go-native-plugin" {
			name, p, err := registry.DownloadGoNativePlugin(c.Config)
			if err != nil {
				return nil, errors.Trace(err)
			}
			registry.RegisterPlugin(registry.FilterPlugin, name, p, true)
		}

		factory, err := registry.GetPlugin(registry.FilterPlugin, c.Type)
		if err != nil {
			return nil, errors.Trace(err)
		}

		filterFactory, ok := factory.(core.IFilterFactory)
		if !ok {
			return nil, errors.Errorf("wrong type: %v", reflect.TypeOf(factory))
		}

		f := filterFactory.NewFilter()

		if err := f.Configure(c.Config); err != nil {
			return nil, errors.Trace(err)
		}

		retFilters = append(retFilters, f)
	}
	return retFilters, nil
}
