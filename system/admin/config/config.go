package config

import (
	"github.com/ponzu-cms/ponzu/management/editor"
	"github.com/ponzu-cms/ponzu/system/item"
)

// Config represents the confirgurable options of the system
type Config struct {
	item.Item
	editor editor.Editor

	Name            string   `json:"name"`
	Domain          string   `json:"domain"`
	HTTPPort        string   `json:"http_port"`
	AdminEmail      string   `json:"admin_email"`
	ClientSecret    string   `json:"client_secret"`
	Etag            string   `json:"etag"`
	CacheInvalidate []string `json:"cache"`
}

// String partially implements item.Identifiable and overrides Item's String()
func (c *Config) String() string { return c.Name }

// Editor partially implements editor.Editable
func (c *Config) Editor() *editor.Editor { return &c.editor }

// MarshalEditor writes a buffer of html to edit a Post and partially implements editor.Editable
func (c *Config) MarshalEditor() ([]byte, error) {
	view, err := editor.Form(c,
		editor.Field{
			View: editor.Input("Name", c, map[string]string{
				"label":       "Site Name",
				"placeholder": "Add a name to this site (internal use only)",
			}),
		},
		editor.Field{
			View: editor.Input("Domain", c, map[string]string{
				"label":       "Domain Name (required for SSL certificate)",
				"placeholder": "e.g. www.example.com or example.com",
			}),
		},
		editor.Field{
			View: editor.Input("HTTPPort", c, map[string]string{
				"type": "hidden",
			}),
		},
		editor.Field{
			View: editor.Input("AdminEmail", c, map[string]string{
				"label": "Adminstrator Email (will be notified of internal system information)",
			}),
		},
		editor.Field{
			View: editor.Input("ClientSecret", c, map[string]string{
				"label":    "Client Secret (used to validate requests, DO NOT SHARE)",
				"disabled": "true",
			}),
		},
		editor.Field{
			View: editor.Input("ClientSecret", c, map[string]string{
				"type": "hidden",
			}),
		},
		editor.Field{
			View: editor.Input("Etag", c, map[string]string{
				"label":    "Etag Header (used for static asset cache)",
				"disabled": "true",
			}),
		},
		editor.Field{
			View: editor.Input("Etag", c, map[string]string{
				"type": "hidden",
			}),
		},
		editor.Field{
			View: editor.Checkbox("CacheInvalidate", c, map[string]string{
				"label": "Invalidate cache on save",
			}, map[string]string{
				"invalidate": "Invalidate Cache",
			}),
		},
	)
	if err != nil {
		return nil, err
	}

	open := []byte(`<div class="card"><form action="/admin/configure" method="post">`)
	close := []byte(`</form></div>`)
	script := []byte(`
	<script>
		$(function() {
			// hide default fields & labels unecessary for the config
			var fields = $('.default-fields');
			fields.css('position', 'relative');
			fields.find('input:not([type=submit])').remove();
			fields.find('label').remove();
			fields.find('button').css({
				position: 'absolute',
				top: '-10px',
				right: '0px'
			});

			var contentOnly = $('.content-only.__ponzu');
			contentOnly.hide();
			contentOnly.find('input, textarea, select').attr('name', '');

			// adjust layout of td so save button is in same location as usual
			fields.find('td').css('float', 'right');

			// stop some fixed config settings from being modified
			fields.find('input[name=client_secret]').attr('name', '');
		});
	</script>
	`)

	view = append(open, view...)
	view = append(view, close...)
	view = append(view, script...)

	return view, nil
}
