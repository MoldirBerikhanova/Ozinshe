package repositories

import (
	"context"
	"errors"
	"goozinshe/models"
)

type MoviesRepository struct {
	db map[int]models.Movie
}

func NewMoviesRepository() *MoviesRepository {
	return &MoviesRepository{
		db: map[int]models.Movie{
			1: {
				Id:          1,
				Title:       "1+1",
				Description: "Пострадав в результате несчастного случая, богатый аристократ Филипп нанимает в помощники человека, который менее всего подходит для этой работы, – молодого жителя предместья Дрисса, только что освободившегося из тюрьмы. Несмотря на то, что Филипп прикован к инвалидному креслу, Дриссу удается привнести в размеренную жизнь аристократа дух приключений.",
				ReleaseYear: 2011,
				Director:    "Оливье Накаш",
				Rating:      0,
				IsWatched:   false,
				TrailerUrl:  "https://www.youtube.com/watch?v=m95M-I7Ij0o&ab_channel=%D0%9A%D0%B8%D0%BD%D0%BE%D0%92%D0%B8%D1%85%D1%80%D1%8C",
				PosterUrl:   "",
				Genres:      make([]models.Genre, 0),
			},
			2: {
				Id:          2,
				Title:       "Интерстеллар",
				Description: "Когда засуха, пыльные бури и вымирание растений приводят человечество к продовольственному кризису, коллектив исследователей и учёных отправляется сквозь червоточину (которая предположительно соединяет области пространства-времени через большое расстояние) в путешествие, чтобы превзойти прежние ограничения для космических путешествий человека и найти планету с подходящими для человечества условиями.",
				ReleaseYear: 2014,
				Director:    "Кристофер Нолан",
				Rating:      0,
				IsWatched:   false,
				TrailerUrl:  "https://www.youtube.com/watch?v=6ybBuTETr3U",
				PosterUrl:   "",
				Genres:      make([]models.Genre, 0),
			},
			3: {
				Id:          3,
				Title:       "Побег из Шоушенка",
				Description: "Бухгалтер Энди Дюфрейн обвинён в убийстве собственной жены и её любовника. Оказавшись в тюрьме под названием Шоушенк, он сталкивается с жестокостью и беззаконием, царящими по обе стороны решётки. Каждый, кто попадает в эти стены, становится их рабом до конца жизни. Но Энди, обладающий живым умом и доброй душой, находит подход как к заключённым, так и к охранникам, добиваясь их особого к себе расположения.",
				ReleaseYear: 1994,
				Director:    "Фрэнк Дарабонт",
				Rating:      0,
				IsWatched:   false,
				TrailerUrl:  "https://www.youtube.com/watch?v=kgAeKpAPOYk&ab_channel=%D0%A2%D1%80%D0%B5%D0%B9%D0%BB%D0%B5%D1%80%D1%8B%D0%BA%D1%84%D0%B8%D0%BB%D1%8C%D0%BC%D0%B0%D0%BC",
				PosterUrl:   "",
				Genres:      make([]models.Genre, 0),
			},
			4: {
				Id:          4,
				Title:       "Бойцовский клуб",
				Description: "Сотрудник страховой компании страдает хронической бессонницей и отчаянно пытается вырваться из мучительно скучной жизни. Однажды в очередной командировке он встречает некоего Тайлера Дёрдена — харизматического торговца мылом с извращенной философией. Тайлер уверен, что самосовершенствование — удел слабых, а единственное, ради чего стоит жить, — саморазрушение. Проходит немного времени, и вот уже новые друзья лупят друг друга почем зря на стоянке перед баром, и очищающий мордобой доставляет им высшее блаженство. Приобщая других мужчин к простым радостям физической жестокости, они основывают тайный Бойцовский клуб, который начинает пользоваться невероятной популярностью.",
				ReleaseYear: 1999,
				Director:    "Дэвид Финчер",
				Rating:      0,
				IsWatched:   false,
				TrailerUrl:  "https://www.youtube.com/watch?v=C7-7qQ61QHU&ab_channel=%D0%A2%D1%80%D0%B5%D0%B9%D0%BB%D0%B5%D1%80%D1%8B%D0%BA%D1%84%D0%B8%D0%BB%D1%8C%D0%BC%D0%B0%D0%BC",
				PosterUrl:   "",
				Genres:      make([]models.Genre, 0),
			},
			5: {
				Id:          5,
				Title:       "Остров проклятых",
				Description: "Два американских судебных пристава отправляются на один из островов в штате Массачусетс, чтобы расследовать исчезновение пациентки клиники для умалишенных преступников. При проведении расследования им придется столкнуться с паутиной лжи, обрушившимся ураганом и смертельным бунтом обитателей клиники.",
				ReleaseYear: 2009,
				Director:    "Мартин Скорсезе",
				Rating:      0,
				IsWatched:   false,
				TrailerUrl:  "https://www.youtube.com/watch?v=_l7R9Rz5URw&ab_channel=%D0%A2%D1%80%D0%B5%D0%B9%D0%BB%D0%B5%D1%80%D1%8B%D0%BA%D1%84%D0%B8%D0%BB%D1%8C%D0%BC%D0%B0%D0%BC",
				PosterUrl:   "",
				Genres:      make([]models.Genre, 0),
			},
		},
	}
}

func (r *MoviesRepository) FindById(c context.Context, id int) (models.Movie, error) {
	movie, ok := r.db[id]
	if !ok {
		return models.Movie{}, errors.New("Movie not found")
	}

	return movie, nil
}

func (r *MoviesRepository) FindAll(c context.Context) []models.Movie {
	movies := make([]models.Movie, 0, len(r.db))
	for _, v := range r.db {
		movies = append(movies, v)
	}

	return movies
}

func (r *MoviesRepository) Create(c context.Context, movie models.Movie) int {
	id := len(r.db) + 1
	movie.Id = id

	r.db[id] = movie

	return id
}

func (r *MoviesRepository) Update(c context.Context, id int, updatedMovie models.Movie) {
	originalMovie := r.db[id]

	originalMovie.Title = updatedMovie.Title
	originalMovie.Description = updatedMovie.Description
	originalMovie.ReleaseYear = updatedMovie.ReleaseYear
	originalMovie.Director = updatedMovie.Director
	originalMovie.TrailerUrl = updatedMovie.TrailerUrl
	originalMovie.Genres = updatedMovie.Genres
	// originalMovie.PosterUrl = updatedMovie.PosterUrl

	r.db[id] = originalMovie
}

func (r *MoviesRepository) Delete(c context.Context, id int) {
	delete(r.db, id)
}
