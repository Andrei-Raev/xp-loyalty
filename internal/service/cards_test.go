package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wintermonth2298/xp-loyalty/internal/model"
	"github.com/wintermonth2298/xp-loyalty/internal/model/mocks"
)

func TestCardsService_GetFormattedCards(t *testing.T) {
	dailyCards := model.Cards{
		{
			Done: 1,
		},
		{
			Done: 0,
		},
		{
			Done: 1,
		},
		{
			Done: 1,
		},
		{
			Done: 0,
		},
		{
			Done: 1,
		},
	}

	constCards := model.Cards{
		{
			Static: model.CardStatic{
				ChainName:  "chain0",
				ChainOrder: 0,
			},
			Done: 1,
		},
		{
			Static: model.CardStatic{
				ChainName:  "chain1",
				ChainOrder: 0,
			},
			Done: 0,
		},
		{
			Static: model.CardStatic{
				ChainName:  "chain1",
				ChainOrder: 1,
			},
			Done: 0,
		},
		{
			Static: model.CardStatic{
				ChainName:  "chain2",
				ChainOrder: 0,
			},
			Done: 1,
		},
		{
			Static: model.CardStatic{
				ChainName:  "chain2",
				ChainOrder: 1,
			},
			Done: 0,
		},
		{
			Static: model.CardStatic{
				ChainName:  "chain2",
				ChainOrder: 2,
			},
			Done: 0,
		},
		{
			Static: model.CardStatic{
				ChainName:  "chain3",
				ChainOrder: 0,
			},
			Done: 1,
		},
		{
			Static: model.CardStatic{
				ChainName:  "chain3",
				ChainOrder: 1,
			},
			Done: 1,
		},
		{
			Static: model.CardStatic{
				ChainName:  "chain3",
				ChainOrder: 2,
			},
			Done: 2,
		},
		{
			Static: model.CardStatic{
				ChainName:  "chain4",
				ChainOrder: 0,
			},
			Done: 1,
		},
		{
			Static: model.CardStatic{
				ChainName:  "chain5",
				ChainOrder: 0,
			},
			Done: 0,
		},
		{
			Static: model.CardStatic{
				ChainName:  "chain6",
				ChainOrder: 0,
			},
			Done: 1,
		},
	}

	selectCards := func(cards model.Cards, indexes []int) model.Cards {
		c := make(model.Cards, 0, len(cards))
		for _, i := range indexes {
			c = append(c, cards[i])
		}
		return c
	}

	wantPending := append(selectCards(dailyCards, []int{1, 4}), selectCards(constCards, []int{0, 1, 4, 8, 9, 10, 11})...)
	wantDone := append(selectCards(dailyCards, []int{0, 2, 3, 5}), selectCards(constCards, []int{0, 3, 6, 7, 8, 9, 11})...)

	cardsRepo := new(mocks.CardsRepositoryMock)
	cardsRepo.On("GetCardsByOwnerPool", mock.Anything, mock.Anything, model.PoolDaily).Return(dailyCards, nil)
	cardsRepo.On("GetCardsByOwnerPool", mock.Anything, mock.Anything, model.PoolConst).Return(constCards, nil)

	t.Run("the only case", func(t *testing.T) {
		s := &CardsService{cardsRepo: cardsRepo}
		gotPending, gotDone, err := s.GetFormattedCards(context.Background(), "user")
		assert.NoError(t, err)

		require.Equal(t, len(wantPending), len(gotPending), "test pending cards len")
		require.Equal(t, len(wantDone), len(gotDone), "test done cards len")

		assert.Equal(t, wantPending, gotPending, "test pending cards")
		assert.Equal(t, wantDone, gotDone, "test done cards")
	})
}

func TestCardsService_Update(t *testing.T) {
	cardsRepo := new(mocks.CardsRepositoryMock)

	type fields struct {
		cardsRepo CardsRepo
	}
	type args struct {
		ctx        context.Context
		id         string
		progress   int
		doneOption float32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		repo   model.Card
		update model.Card
	}{
		{
			name: "case ordinary 1",
			fields: fields{
				cardsRepo: cardsRepo,
			},
			args: args{
				progress:   0,
				doneOption: 0,
			},
			want: 100,
			repo: model.Card{
				Static: model.CardStatic{
					Type: "ordinary",
					OrdSettings: &model.OrdSettings{
						Award: model.Award{
							XPoints: 100,
						},
					},
				},
				Done:     0,
				Progress: 0,
				History:  []int{},
			},
			update: model.Card{
				Static: model.CardStatic{
					Type: "ordinary",
					OrdSettings: &model.OrdSettings{
						Award: model.Award{
							XPoints: 100,
						},
					},
				},
				Done:     1,
				Progress: 0,
				History:  []int{},
			},
		},

		{
			name: "case ordinary 2",
			fields: fields{
				cardsRepo: cardsRepo,
			},
			args: args{
				progress:   100,
				doneOption: 100,
			},
			want: 100,
			repo: model.Card{
				Static: model.CardStatic{
					Type: model.TypeOrdinary,
					OrdSettings: &model.OrdSettings{
						Award: model.Award{
							XPoints: 100,
						},
					},
				},
				Done:     3,
				Progress: 0,
				History:  []int{},
			},
			update: model.Card{
				Static: model.CardStatic{
					Type: model.TypeOrdinary,
					OrdSettings: &model.OrdSettings{
						Award: model.Award{
							XPoints: 100,
						},
					},
				},
				Done:     4,
				Progress: 0,
				History:  []int{},
			},
		},

		{
			name: "case progerss 1",
			fields: fields{
				cardsRepo: cardsRepo,
			},
			args: args{
				progress:   5,
				doneOption: 0,
			},
			want: 0,
			repo: model.Card{
				Static: model.CardStatic{
					Type: model.TypeProgress,
					PrgSettings: &model.PrgSettings{
						Award: model.Award{
							XPoints: 200,
						},
						MaxProgress: 20,
					},
				},
				Done:     0,
				Progress: 10,
				History:  []int{},
			},
			update: model.Card{
				Static: model.CardStatic{
					Type: model.TypeProgress,
					PrgSettings: &model.PrgSettings{
						Award: model.Award{
							XPoints: 200,
						},
						MaxProgress: 20,
					},
				},
				Done:     0,
				Progress: 15,
				History:  []int{},
			},
		},

		{
			name: "case progerss 2",
			fields: fields{
				cardsRepo: cardsRepo,
			},
			args: args{
				progress:   100,
				doneOption: 0,
			},
			want: 200,
			repo: model.Card{
				Static: model.CardStatic{
					Type: model.TypeProgress,
					PrgSettings: &model.PrgSettings{
						Award: model.Award{
							XPoints: 200,
						},
						MaxProgress: 20,
					},
				},
				Done:     0,
				Progress: 10,
				History:  []int{},
			},
			update: model.Card{
				Static: model.CardStatic{
					Type: model.TypeProgress,
					PrgSettings: &model.PrgSettings{
						Award: model.Award{
							XPoints: 200,
						},
						MaxProgress: 20,
					},
				},
				Done:     1,
				Progress: 0,
				History:  []int{},
			},
		},

		{
			name: "case options 1",
			fields: fields{
				cardsRepo: cardsRepo,
			},
			args: args{
				progress:   100,
				doneOption: 5,
			},
			want: 0,
			repo: model.Card{
				Static: model.CardStatic{
					Type: model.TypeOptions,
					OptSettings: &model.OptSettings{
						Awards: []model.Award{{
							XPoints: 100,
						}, {
							XPoints: 200,
						}, {
							XPoints: 300,
						}},
						Options: []float32{10, 20, 30},
					},
				},
				Done:     0,
				Progress: 0,
				History:  []int{},
			},
			update: model.Card{
				Static: model.CardStatic{
					Type: model.TypeOptions,
					OptSettings: &model.OptSettings{
						Awards: []model.Award{{
							XPoints: 100,
						}, {
							XPoints: 200,
						}, {
							XPoints: 300,
						}},
						Options: []float32{10, 20, 30},
					},
				},
				Done:     0,
				Progress: 0,
				History:  []int{},
			},
		},

		{
			name: "case options 2",
			fields: fields{
				cardsRepo: cardsRepo,
			},
			args: args{
				progress:   100,
				doneOption: 10,
			},
			want: 100,
			repo: model.Card{
				Static: model.CardStatic{
					Type: model.TypeOptions,
					OptSettings: &model.OptSettings{
						Awards: []model.Award{{
							XPoints: 100,
						}, {
							XPoints: 200,
						}, {
							XPoints: 300,
						}},
						Options: []float32{10, 20, 30},
					},
				},
				Done:     0,
				Progress: 0,
				History:  []int{},
			},
			update: model.Card{
				Static: model.CardStatic{
					Type: model.TypeOptions,
					OptSettings: &model.OptSettings{
						Awards: []model.Award{{
							XPoints: 100,
						}, {
							XPoints: 200,
						}, {
							XPoints: 300,
						}},
						Options: []float32{10, 20, 30},
					},
				},
				Done:     1,
				Progress: 0,
				History:  []int{10},
			},
		},

		{
			name: "case options 3",
			fields: fields{
				cardsRepo: cardsRepo,
			},
			args: args{
				doneOption: 11,
			},
			want: 100,
			repo: model.Card{
				Static: model.CardStatic{
					Type: model.TypeOptions,
					OptSettings: &model.OptSettings{
						Awards: []model.Award{{
							XPoints: 100,
						}, {
							XPoints: 200,
						}, {
							XPoints: 300,
						}},
						Options: []float32{10, 20, 30},
					},
				},
				Done:    2,
				History: []int{},
			},
			update: model.Card{
				Static: model.CardStatic{
					Type: model.TypeOptions,
					OptSettings: &model.OptSettings{
						Awards: []model.Award{{
							XPoints: 100,
						}, {
							XPoints: 200,
						}, {
							XPoints: 300,
						}},
						Options: []float32{10, 20, 30},
					},
				},
				Done:    3,
				History: []int{11},
			},
		},

		{
			name: "case options 4",
			fields: fields{
				cardsRepo: cardsRepo,
			},
			args: args{
				doneOption: 20,
			},
			want: 200,
			repo: model.Card{
				Static: model.CardStatic{
					Type: model.TypeOptions,
					OptSettings: &model.OptSettings{
						Awards: []model.Award{{
							XPoints: 100,
						}, {
							XPoints: 200,
						}, {
							XPoints: 300,
						}},
						Options: []float32{10, 20, 30},
					},
				},
				Done:    2,
				History: []int{30},
			},
			update: model.Card{
				Static: model.CardStatic{
					Type: model.TypeOptions,
					OptSettings: &model.OptSettings{
						Awards: []model.Award{{
							XPoints: 100,
						}, {
							XPoints: 200,
						}, {
							XPoints: 300,
						}},
						Options: []float32{10, 20, 30},
					},
				},
				Done:    3,
				History: []int{30, 20},
			},
		},

		{
			name: "case options 5",
			fields: fields{
				cardsRepo: cardsRepo,
			},
			args: args{
				doneOption: 30,
			},
			want: 300,
			repo: model.Card{
				Static: model.CardStatic{
					Type: model.TypeOptions,
					OptSettings: &model.OptSettings{
						Awards: []model.Award{{
							XPoints: 100,
						}, {
							XPoints: 200,
						}, {
							XPoints: 300,
						}},
						Options: []float32{10, 20, 30},
					},
				},
				Done:    2,
				History: []int{30},
			},
			update: model.Card{
				Static: model.CardStatic{
					Type: model.TypeOptions,
					OptSettings: &model.OptSettings{
						Awards: []model.Award{{
							XPoints: 100,
						}, {
							XPoints: 200,
						}, {
							XPoints: 300,
						}},
						Options: []float32{10, 20, 30},
					},
				},
				Done:    3,
				History: []int{30, 30},
			},
		},

		{
			name: "case options 6",
			fields: fields{
				cardsRepo: cardsRepo,
			},
			args: args{
				doneOption: 100,
			},
			want: 300,
			repo: model.Card{
				Static: model.CardStatic{
					Type: model.TypeOptions,
					OptSettings: &model.OptSettings{
						Awards: []model.Award{{
							XPoints: 100,
						}, {
							XPoints: 200,
						}, {
							XPoints: 300,
						}},
						Options: []float32{10, 20, 30},
					},
				},
				Done:    2,
				History: []int{30},
			},
			update: model.Card{
				Static: model.CardStatic{
					Type: model.TypeOptions,
					OptSettings: &model.OptSettings{
						Awards: []model.Award{{
							XPoints: 100,
						}, {
							XPoints: 200,
						}, {
							XPoints: 300,
						}},
						Options: []float32{10, 20, 30},
					},
				},
				Done:    3,
				History: []int{30, 100},
			},
		},
	}

	for _, tt := range tests {
		call1 := cardsRepo.On("Get", mock.Anything, tt.args.id).Return(tt.repo, nil).Once()
		call2 := cardsRepo.On("Update", mock.Anything, tt.update).Return(nil).Maybe()
		s := &CardsService{cardsRepo: tt.fields.cardsRepo}
		_, got, err := s.Update(tt.args.ctx, tt.args.id, tt.args.progress, tt.args.doneOption)

		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, err)
			require.Equal(t, tt.want, got)

			cardsRepo.AssertExpectations(t)
			call1.Unset()
			call2.Unset()
		})
	}
}

func TestCardsService_UpdateConstCards(t *testing.T) {
	type args struct {
		ctx   context.Context
		users []model.User
	}
	test := struct {
		name            string
		args            args
		repoCardsUser1  model.Cards
		repoCardsUser2  model.Cards
		repoCardsUser3  model.Cards
		repoStaticCards model.CardsStatic
		wantUser1       model.Cards
		wantUser2       model.Cards
		wantUser3       model.Cards
	}{
		name: "case1",
		args: args{
			ctx: context.Background(),
			users: []model.User{{
				CredentialsSecure: model.CredentialsSecure{
					Username: "user1",
				},
			}, {
				CredentialsSecure: model.CredentialsSecure{
					Username: "user2",
				},
			}, {
				CredentialsSecure: model.CredentialsSecure{
					Username: "user3",
				},
			}},
		},
		repoCardsUser1: model.Cards{{
			OwnerUsername: "user1",
			Static: model.CardStatic{
				ID: "1",
			},
		}, {
			OwnerUsername: "user1",
			Static: model.CardStatic{
				ID: "10",
			},
		}},

		repoCardsUser2: model.Cards{{
			OwnerUsername: "user2",
			Static: model.CardStatic{
				ID: "1",
			},
		}},

		repoCardsUser3: model.Cards{},

		repoStaticCards: model.CardsStatic{{
			ID: "1",
		}, {
			ID: "10",
		}},

		wantUser1: model.Cards{},
		wantUser2: model.Cards{{
			OwnerUsername: "user2",
			Static: model.CardStatic{
				ID: "10",
			},
		}},

		wantUser3: model.Cards{{
			OwnerUsername: "user3",
			Static: model.CardStatic{
				ID: "1",
			},
		}, {
			OwnerUsername: "user3",
			Static: model.CardStatic{
				ID: "10",
			},
		}},
	}

	cardsRepo := new(mocks.CardsRepositoryMock)
	cardsRepo.On("GetStaticByPool", mock.Anything, mock.Anything).Return(test.repoStaticCards, nil).Once()
	cardsRepo.On("GetCardsByOwner", mock.Anything, "user1").Return(test.repoCardsUser1, nil)
	cardsRepo.On("GetCardsByOwner", mock.Anything, "user2").Return(test.repoCardsUser2, nil)
	cardsRepo.On("GetCardsByOwner", mock.Anything, "user3").Return(test.repoCardsUser3, nil)
	cardsRepo.On("Create", mock.Anything, test.wantUser2[0]).Return(nil)
	cardsRepo.On("Create", mock.Anything, test.wantUser3[0]).Return(nil)
	cardsRepo.On("Create", mock.Anything, test.wantUser3[1]).Return(nil)
	s := &CardsService{cardsRepo: cardsRepo}

	t.Run(test.name, func(t *testing.T) {
		err := s.UpdateConstCards(test.args.ctx, test.args.users)
		require.NoError(t, err)
		cardsRepo.AssertExpectations(t)
	})
}
