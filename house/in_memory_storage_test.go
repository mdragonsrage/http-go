package house

import "testing"

func TestInMemoryStorage_create(t *testing.T) {
	type fields struct {
		data map[string]bool
	}
	type args struct {
		lb Lightbulb
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantCount int
	}{
		{
			name: "empty_data_should:_have_one_entry_after_create",
			fields: fields{
				data: map[string]bool{},
			},
			args: args{
				lb: Lightbulb{
					Name: "livingroom",
				},
			},
			wantErr:   false,
			wantCount: 1,
		}, {
			name: "if_one_item_exists_data_should_have_two_entries_after_create",
			fields: fields{
				data: map[string]bool{
					"bedroom": false,
				},
			},
			args: args{
				lb: Lightbulb{
					Name: "living-room",
				},
			},
			wantErr:   false,
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &InMemoryStorage{
				data: tt.fields.data,
			}
			if err := db.Create(tt.args.lb); (err != nil) != tt.wantErr {
				t.Errorf("InMemoryStorage.Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			items, _ := db.GetAll()
			itemCount := len(items)

			if itemCount != tt.wantCount {
				t.Errorf("itemCount %d, wantCount %d", itemCount, tt.wantCount)
			}
		})
	}

}
