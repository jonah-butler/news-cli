# Read The News

A simple Command Line utility for reading news results from the command line

## Pre-requisites

The following env fields are requires to run

- BASE\*URL: ex. "https://somenewsagency.com"

  _The following var syntax utilizes CSS selectors for interacting with various pages of the site_

- SEARCH_URL: ex. "https://somenewsagency.com/search/?some=defaults&query=settings"
- ARTICLE_BODY = "section#main-page-container"
- ARTICLE_TITLE = "h1.headline span"
- ARTICLE_TEXT = ".article-text:not(.trinity-skip-it, .html-content > .lee-article-text)"
- RESULTS_CONTAINER = "div#results-col"
- RESULTS_LINK = "h3.tnt-headline a"

## Disclaimer

Use at your own risk and be confident of your access rights when scraping news content
