func RequestNijisanjiWikiRawSchedule() (*goquery.Document, error) {
	res, err := http.Get("https://wikiwiki.jp/nijisanji/%E9%85%8D%E4%BF%A1%E4%BA%88%E5%AE%9A%E3%83%AA%E3%82%B9%E3%83%88")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
