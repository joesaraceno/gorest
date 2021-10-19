package main

import "errors"

type Article struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

var Articles []Article = []Article{
	{Id: 1, Title: "Hello", Desc: "Article Description", Content: "Article Content"},
	{Id: 2, Title: "Hello2", Desc: "Article Description", Content: "Article Content"},
}

func (a *Article) getArticle() error {
	for _, article := range Articles {
		if article.Id == a.Id {
			a.Desc = article.Desc
			a.Title = article.Title
			a.Content = article.Content
			return nil
		}
	}
	return errors.New("not found")
}

func getArticlesFromRepo() ([]Article, error) {
	return Articles, nil
}

func (a *Article) createArticle() error {
	Articles = append(Articles, *a)
	return nil
}

func (a *Article) deleteArticle() error {
	for i, article := range Articles {
		if article.Id == a.Id {
			Articles = append(Articles[:i], Articles[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}
