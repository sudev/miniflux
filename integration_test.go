// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// +build integration

package main

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/miniflux/miniflux-go"
)

const (
	testBaseURL          = "http://127.0.0.1:8080/"
	testAdminUsername    = "admin"
	testAdminPassword    = "test123"
	testStandardPassword = "secret"
	testFeedURL          = "https://github.com/miniflux/miniflux/commits/master.atom"
	testFeedTitle        = "Recent Commits to miniflux:master"
	testWebsiteURL       = "https://github.com/miniflux/miniflux/commits/master"
)

func TestWithBadEndpoint(t *testing.T) {
	client := miniflux.NewClient("bad url", testAdminUsername, testAdminPassword)
	_, err := client.Users()
	if err == nil {
		t.Fatal(`Using a bad url should raise an error`)
	}
}

func TestWithWrongCredentials(t *testing.T) {
	client := miniflux.NewClient(testBaseURL, "invalid", "invalid")
	_, err := client.Users()
	if err == nil {
		t.Fatal(`Using bad credentials should raise an error`)
	}

	if err != miniflux.ErrNotAuthorized {
		t.Fatal(`A "Not Authorized" error should be raised`)
	}
}

func TestGetCurrentLoggedUser(t *testing.T) {
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	user, err := client.Me()
	if err != nil {
		t.Fatal(err)
	}

	if user.ID == 0 {
		t.Fatalf(`Invalid userID, got %q`, user.ID)
	}

	if user.Username != testAdminUsername {
		t.Fatalf(`Invalid username, got %q`, user.Username)
	}
}

func TestGetUsers(t *testing.T) {
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	users, err := client.Users()
	if err != nil {
		t.Fatal(err)
	}

	if len(users) == 0 {
		t.Fatal("The list of users is empty")
	}

	if users[0].ID == 0 {
		t.Fatalf(`Invalid userID, got "%v"`, users[0].ID)
	}

	if users[0].Username != testAdminUsername {
		t.Fatalf(`Invalid username, got "%v" instead of "%v"`, users[0].Username, testAdminUsername)
	}

	if users[0].Password != "" {
		t.Fatalf(`Invalid password, got "%v"`, users[0].Password)
	}

	if users[0].Language != "en_US" {
		t.Fatalf(`Invalid language, got "%v"`, users[0].Language)
	}

	if users[0].Theme != "default" {
		t.Fatalf(`Invalid theme, got "%v"`, users[0].Theme)
	}

	if users[0].Timezone != "UTC" {
		t.Fatalf(`Invalid timezone, got "%v"`, users[0].Timezone)
	}

	if !users[0].IsAdmin {
		t.Fatalf(`Invalid role, got "%v"`, users[0].IsAdmin)
	}
}

func TestCreateStandardUser(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	user, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID == 0 {
		t.Fatalf(`Invalid userID, got "%v"`, user.ID)
	}

	if user.Username != username {
		t.Fatalf(`Invalid username, got "%v" instead of "%v"`, user.Username, username)
	}

	if user.Password != "" {
		t.Fatalf(`Invalid password, got "%v"`, user.Password)
	}

	if user.Language != "en_US" {
		t.Fatalf(`Invalid language, got "%v"`, user.Language)
	}

	if user.Theme != "default" {
		t.Fatalf(`Invalid theme, got "%v"`, user.Theme)
	}

	if user.Timezone != "UTC" {
		t.Fatalf(`Invalid timezone, got "%v"`, user.Timezone)
	}

	if user.IsAdmin {
		t.Fatalf(`Invalid role, got "%v"`, user.IsAdmin)
	}

	if user.LastLoginAt != nil {
		t.Fatalf(`Invalid last login date, got "%v"`, user.LastLoginAt)
	}
}

func TestRemoveUser(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	user, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	if err := client.DeleteUser(user.ID); err != nil {
		t.Fatalf(`Unable to remove user: "%v"`, err)
	}
}

func TestGetUserByID(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	user, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.UserByID(99999)
	if err == nil {
		t.Fatal(`Should returns a 404`)
	}

	user, err = client.UserByID(user.ID)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID == 0 {
		t.Fatalf(`Invalid userID, got "%v"`, user.ID)
	}

	if user.Username != username {
		t.Fatalf(`Invalid username, got "%v" instead of "%v"`, user.Username, username)
	}

	if user.Password != "" {
		t.Fatalf(`Invalid password, got "%v"`, user.Password)
	}

	if user.Language != "en_US" {
		t.Fatalf(`Invalid language, got "%v"`, user.Language)
	}

	if user.Theme != "default" {
		t.Fatalf(`Invalid theme, got "%v"`, user.Theme)
	}

	if user.Timezone != "UTC" {
		t.Fatalf(`Invalid timezone, got "%v"`, user.Timezone)
	}

	if user.IsAdmin {
		t.Fatalf(`Invalid role, got "%v"`, user.IsAdmin)
	}

	if user.LastLoginAt != nil {
		t.Fatalf(`Invalid last login date, got "%v"`, user.LastLoginAt)
	}
}

func TestGetUserByUsername(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	user, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.UserByUsername("missinguser")
	if err == nil {
		t.Fatal(`Should returns a 404`)
	}

	user, err = client.UserByUsername(username)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID == 0 {
		t.Fatalf(`Invalid userID, got "%v"`, user.ID)
	}

	if user.Username != username {
		t.Fatalf(`Invalid username, got "%v" instead of "%v"`, user.Username, username)
	}

	if user.Password != "" {
		t.Fatalf(`Invalid password, got "%v"`, user.Password)
	}

	if user.Language != "en_US" {
		t.Fatalf(`Invalid language, got "%v"`, user.Language)
	}

	if user.Theme != "default" {
		t.Fatalf(`Invalid theme, got "%v"`, user.Theme)
	}

	if user.Timezone != "UTC" {
		t.Fatalf(`Invalid timezone, got "%v"`, user.Timezone)
	}

	if user.IsAdmin {
		t.Fatalf(`Invalid role, got "%v"`, user.IsAdmin)
	}

	if user.LastLoginAt != nil {
		t.Fatalf(`Invalid last login date, got "%v"`, user.LastLoginAt)
	}
}

func TestUpdateUserTheme(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	user, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	theme := "black"
	user, err = client.UpdateUser(user.ID, &miniflux.UserModification{Theme: &theme})
	if err != nil {
		t.Fatal(err)
	}

	if user.Theme != theme {
		t.Fatalf(`Unable to update user Theme: got "%v" instead of "%v"`, user.Theme, theme)
	}
}

func TestUpdateUserThemeWithInvalidValue(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	user, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	theme := "something that doesn't exists"
	_, err = client.UpdateUser(user.ID, &miniflux.UserModification{Theme: &theme})
	if err == nil {
		t.Fatal(`Updating a user Theme with an invalid value should raise an error`)
	}
}

func TestCannotCreateDuplicateUser(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	_, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.CreateUser(username, testStandardPassword, false)
	if err == nil {
		t.Fatal(`Duplicate users should not be allowed`)
	}
}

func TestCannotListUsersAsNonAdmin(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	_, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	client = miniflux.NewClient(testBaseURL, username, testStandardPassword)
	_, err = client.Users()
	if err == nil {
		t.Fatal(`Standard users should not be able to list any users`)
	}

	if err != miniflux.ErrForbidden {
		t.Fatal(`A "Forbidden" error should be raised`)
	}
}

func TestCannotGetUserAsNonAdmin(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	user, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	client = miniflux.NewClient(testBaseURL, username, testStandardPassword)
	_, err = client.UserByID(user.ID)
	if err == nil {
		t.Fatal(`Standard users should not be able to get any users`)
	}

	if err != miniflux.ErrForbidden {
		t.Fatal(`A "Forbidden" error should be raised`)
	}
}

func TestCannotUpdateUserAsNonAdmin(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	user, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	client = miniflux.NewClient(testBaseURL, username, testStandardPassword)
	_, err = client.UpdateUser(user.ID, &miniflux.UserModification{})
	if err == nil {
		t.Fatal(`Standard users should not be able to update any users`)
	}

	if err != miniflux.ErrForbidden {
		t.Fatal(`A "Forbidden" error should be raised`)
	}
}

func TestCannotCreateUserAsNonAdmin(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	_, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	client = miniflux.NewClient(testBaseURL, username, testStandardPassword)
	_, err = client.CreateUser(username, testStandardPassword, false)
	if err == nil {
		t.Fatal(`Standard users should not be able to create users`)
	}

	if err != miniflux.ErrForbidden {
		t.Fatal(`A "Forbidden" error should be raised`)
	}
}

func TestCannotDeleteUserAsNonAdmin(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	user, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	client = miniflux.NewClient(testBaseURL, username, testStandardPassword)
	err = client.DeleteUser(user.ID)
	if err == nil {
		t.Fatal(`Standard users should not be able to remove any users`)
	}

	if err != miniflux.ErrForbidden {
		t.Fatal(`A "Forbidden" error should be raised`)
	}
}

func TestCreateCategory(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	user, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	categoryName := "My category"
	client = miniflux.NewClient(testBaseURL, username, testStandardPassword)
	category, err := client.CreateCategory(categoryName)
	if err != nil {
		t.Fatal(err)
	}

	if category.ID == 0 {
		t.Fatalf(`Invalid categoryID, got "%v"`, category.ID)
	}

	if category.UserID != user.ID {
		t.Fatalf(`Invalid userID, got "%v" instead of "%v"`, category.UserID, user.ID)
	}

	if category.Title != categoryName {
		t.Fatalf(`Invalid title, got "%v" instead of "%v"`, category.Title, categoryName)
	}
}

func TestCreateCategoryWithEmptyTitle(t *testing.T) {
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	_, err := client.CreateCategory("")
	if err == nil {
		t.Fatal(`The category title should be mandatory`)
	}
}

func TestCannotCreateDuplicatedCategory(t *testing.T) {
	client := createClient(t)

	categoryName := "My category"
	_, err := client.CreateCategory(categoryName)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.CreateCategory(categoryName)
	if err == nil {
		t.Fatal(`Duplicated categories should not be allowed`)
	}
}

func TestUpdateCategory(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	user, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	categoryName := "My category"
	client = miniflux.NewClient(testBaseURL, username, testStandardPassword)
	category, err := client.CreateCategory(categoryName)
	if err != nil {
		t.Fatal(err)
	}

	categoryName = "Updated category"
	category, err = client.UpdateCategory(category.ID, categoryName)
	if err != nil {
		t.Fatal(err)
	}

	if category.ID == 0 {
		t.Fatalf(`Invalid categoryID, got "%v"`, category.ID)
	}

	if category.UserID != user.ID {
		t.Fatalf(`Invalid userID, got "%v" instead of "%v"`, category.UserID, user.ID)
	}

	if category.Title != categoryName {
		t.Fatalf(`Invalid title, got "%v" instead of "%v"`, category.Title, categoryName)
	}
}

func TestListCategories(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	user, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	categoryName := "My category"
	client = miniflux.NewClient(testBaseURL, username, testStandardPassword)
	_, err = client.CreateCategory(categoryName)
	if err != nil {
		t.Fatal(err)
	}

	categories, err := client.Categories()
	if err != nil {
		t.Fatal(err)
	}

	if len(categories) != 2 {
		t.Fatalf(`Invalid number of categories, got "%v" instead of "%v"`, len(categories), 2)
	}

	if categories[0].ID == 0 {
		t.Fatalf(`Invalid categoryID, got "%v"`, categories[0].ID)
	}

	if categories[0].UserID != user.ID {
		t.Fatalf(`Invalid userID, got "%v" instead of "%v"`, categories[0].UserID, user.ID)
	}

	if categories[0].Title != "All" {
		t.Fatalf(`Invalid title, got "%v" instead of "%v"`, categories[0].Title, "All")
	}

	if categories[1].ID == 0 {
		t.Fatalf(`Invalid categoryID, got "%v"`, categories[0].ID)
	}

	if categories[1].UserID != user.ID {
		t.Fatalf(`Invalid userID, got "%v" instead of "%v"`, categories[1].UserID, user.ID)
	}

	if categories[1].Title != categoryName {
		t.Fatalf(`Invalid title, got "%v" instead of "%v"`, categories[1].Title, categoryName)
	}
}

func TestDeleteCategory(t *testing.T) {
	client := createClient(t)

	category, err := client.CreateCategory("My category")
	if err != nil {
		t.Fatal(err)
	}

	err = client.DeleteCategory(category.ID)
	if err != nil {
		t.Fatal(`Removing a category should not raise any error`)
	}
}

func TestCannotDeleteCategoryOfAnotherUser(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	categories, err := client.Categories()
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	client = miniflux.NewClient(testBaseURL, username, testStandardPassword)
	err = client.DeleteCategory(categories[0].ID)
	if err == nil {
		t.Fatal(`Removing a category that belongs to another user should be forbidden`)
	}
}

func TestDiscoverSubscriptions(t *testing.T) {
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	subscriptions, err := client.Discover(testWebsiteURL)
	if err != nil {
		t.Fatal(err)
	}

	if len(subscriptions) != 1 {
		t.Fatalf(`Invalid number of subscriptions, got "%v" instead of "%v"`, len(subscriptions), 2)
	}

	if subscriptions[0].Title != testFeedTitle {
		t.Fatalf(`Invalid feed title, got "%v" instead of "%v"`, subscriptions[0].Title, testFeedTitle)
	}

	if subscriptions[0].Type != "atom" {
		t.Fatalf(`Invalid feed type, got "%v" instead of "%v"`, subscriptions[0].Type, "atom")
	}

	if subscriptions[0].URL != testFeedURL {
		t.Fatalf(`Invalid feed URL, got "%v" instead of "%v"`, subscriptions[0].URL, testFeedURL)
	}
}

func TestCreateFeed(t *testing.T) {
	client := createClient(t)
	feed, _ := createFeed(t, client)

	if feed.ID == 0 {
		t.Fatalf(`Invalid feed ID, got %q`, feed.ID)
	}
}

func TestCannotCreateDuplicatedFeed(t *testing.T) {
	client := createClient(t)
	feed, category := createFeed(t, client)

	_, err := client.CreateFeed(feed.FeedURL, category.ID)
	if err == nil {
		t.Fatal(`Duplicated feeds should not be allowed`)
	}
}

func TestCreateFeedWithInexistingCategory(t *testing.T) {
	client := createClient(t)

	_, err := client.CreateFeed(testFeedURL, -1)
	if err == nil {
		t.Fatal(`Feeds should not be created with inexisting category`)
	}
}

func TestExport(t *testing.T) {
	client := createClient(t)

	output, err := client.Export()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(string(output), "<?xml") {
		t.Fatalf(`Invalid OPML export, got "%s"`, string(output))
	}
}

func TestImport(t *testing.T) {
	client := createClient(t)

	data := `<?xml version="1.0" encoding="UTF-8"?>
    <opml version="2.0">
        <body>
            <outline text="Test Category">
				<outline title="Test" text="Test" xmlUrl="` + testFeedURL + `" htmlUrl="` + testWebsiteURL + `"></outline>
			</outline>
		</body>
	</opml>`

	b := bytes.NewReader([]byte(data))
	err := client.Import(ioutil.NopCloser(b))
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateFeedURL(t *testing.T) {
	client := createClient(t)
	feed, _ := createFeed(t, client)

	url := "test"
	updatedFeed, err := client.UpdateFeed(feed.ID, &miniflux.FeedModification{FeedURL: &url})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.FeedURL != url {
		t.Fatalf(`Wrong FeedURL, got %q instead of %q`, updatedFeed.FeedURL, url)
	}

	url = ""
	updatedFeed, err = client.UpdateFeed(feed.ID, &miniflux.FeedModification{FeedURL: &url})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.FeedURL == "" {
		t.Fatalf(`The FeedURL should not be empty`)
	}
}

func TestUpdateFeedSiteURL(t *testing.T) {
	client := createClient(t)
	feed, _ := createFeed(t, client)

	url := "test"
	updatedFeed, err := client.UpdateFeed(feed.ID, &miniflux.FeedModification{SiteURL: &url})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.SiteURL != url {
		t.Fatalf(`Wrong SiteURL, got %q instead of %q`, updatedFeed.SiteURL, url)
	}

	url = ""
	updatedFeed, err = client.UpdateFeed(feed.ID, &miniflux.FeedModification{SiteURL: &url})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.SiteURL == "" {
		t.Fatalf(`The SiteURL should not be empty`)
	}
}

func TestUpdateFeedTitle(t *testing.T) {
	client := createClient(t)
	feed, _ := createFeed(t, client)

	newTitle := "My new feed"
	updatedFeed, err := client.UpdateFeed(feed.ID, &miniflux.FeedModification{Title: &newTitle})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.Title != newTitle {
		t.Fatalf(`Wrong title, got %q instead of %q`, updatedFeed.Title, newTitle)
	}

	newTitle = ""
	updatedFeed, err = client.UpdateFeed(feed.ID, &miniflux.FeedModification{Title: &newTitle})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.Title == "" {
		t.Fatalf(`The Title should not be empty`)
	}
}

func TestUpdateFeedCrawler(t *testing.T) {
	client := createClient(t)
	feed, _ := createFeed(t, client)

	crawler := true
	updatedFeed, err := client.UpdateFeed(feed.ID, &miniflux.FeedModification{Crawler: &crawler})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.Crawler != crawler {
		t.Fatalf(`Wrong crawler value, got "%v" instead of "%v"`, updatedFeed.Crawler, crawler)
	}

	if updatedFeed.Title != feed.Title {
		t.Fatalf(`The titles should be the same after update`)
	}

	crawler = false
	updatedFeed, err = client.UpdateFeed(feed.ID, &miniflux.FeedModification{Crawler: &crawler})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.Crawler != crawler {
		t.Fatalf(`Wrong crawler value, got "%v" instead of "%v"`, updatedFeed.Crawler, crawler)
	}
}

func TestUpdateFeedScraperRules(t *testing.T) {
	client := createClient(t)
	feed, _ := createFeed(t, client)

	scraperRules := "test"
	updatedFeed, err := client.UpdateFeed(feed.ID, &miniflux.FeedModification{ScraperRules: &scraperRules})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.ScraperRules != scraperRules {
		t.Fatalf(`Wrong ScraperRules value, got "%v" instead of "%v"`, updatedFeed.ScraperRules, scraperRules)
	}

	scraperRules = ""
	updatedFeed, err = client.UpdateFeed(feed.ID, &miniflux.FeedModification{ScraperRules: &scraperRules})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.ScraperRules != scraperRules {
		t.Fatalf(`Wrong ScraperRules value, got "%v" instead of "%v"`, updatedFeed.ScraperRules, scraperRules)
	}
}

func TestUpdateFeedRewriteRules(t *testing.T) {
	client := createClient(t)
	feed, _ := createFeed(t, client)

	rewriteRules := "test"
	updatedFeed, err := client.UpdateFeed(feed.ID, &miniflux.FeedModification{RewriteRules: &rewriteRules})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.RewriteRules != rewriteRules {
		t.Fatalf(`Wrong RewriteRules value, got "%v" instead of "%v"`, updatedFeed.RewriteRules, rewriteRules)
	}

	rewriteRules = ""
	updatedFeed, err = client.UpdateFeed(feed.ID, &miniflux.FeedModification{RewriteRules: &rewriteRules})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.RewriteRules != rewriteRules {
		t.Fatalf(`Wrong RewriteRules value, got "%v" instead of "%v"`, updatedFeed.RewriteRules, rewriteRules)
	}
}

func TestUpdateFeedUsername(t *testing.T) {
	client := createClient(t)
	feed, _ := createFeed(t, client)

	username := "test"
	updatedFeed, err := client.UpdateFeed(feed.ID, &miniflux.FeedModification{Username: &username})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.Username != username {
		t.Fatalf(`Wrong Username value, got "%v" instead of "%v"`, updatedFeed.Username, username)
	}

	username = ""
	updatedFeed, err = client.UpdateFeed(feed.ID, &miniflux.FeedModification{Username: &username})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.Username != username {
		t.Fatalf(`Wrong Username value, got "%v" instead of "%v"`, updatedFeed.Username, username)
	}
}

func TestUpdateFeedPassword(t *testing.T) {
	client := createClient(t)
	feed, _ := createFeed(t, client)

	password := "test"
	updatedFeed, err := client.UpdateFeed(feed.ID, &miniflux.FeedModification{Password: &password})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.Password != password {
		t.Fatalf(`Wrong Password value, got "%v" instead of "%v"`, updatedFeed.Password, password)
	}

	password = ""
	updatedFeed, err = client.UpdateFeed(feed.ID, &miniflux.FeedModification{Password: &password})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.Password != password {
		t.Fatalf(`Wrong Password value, got "%v" instead of "%v"`, updatedFeed.Password, password)
	}
}

func TestUpdateFeedCategory(t *testing.T) {
	client := createClient(t)
	feed, _ := createFeed(t, client)

	newCategory, err := client.CreateCategory("my new category")
	if err != nil {
		t.Fatal(err)
	}

	updatedFeed, err := client.UpdateFeed(feed.ID, &miniflux.FeedModification{CategoryID: &newCategory.ID})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.Category.ID != newCategory.ID {
		t.Fatalf(`Wrong CategoryID value, got "%v" instead of "%v"`, updatedFeed.Category.ID, newCategory.ID)
	}

	categoryID := int64(0)
	updatedFeed, err = client.UpdateFeed(feed.ID, &miniflux.FeedModification{CategoryID: &categoryID})
	if err != nil {
		t.Fatal(err)
	}

	if updatedFeed.Category.ID == 0 {
		t.Fatalf(`The CategoryID must defined`)
	}
}

func TestDeleteFeed(t *testing.T) {
	client := createClient(t)
	feed, _ := createFeed(t, client)
	if err := client.DeleteFeed(feed.ID); err != nil {
		t.Fatal(err)
	}
}

func TestRefreshFeed(t *testing.T) {
	client := createClient(t)
	feed, _ := createFeed(t, client)
	if err := client.RefreshFeed(feed.ID); err != nil {
		t.Fatal(err)
	}
}

func TestGetFeed(t *testing.T) {
	client := createClient(t)
	feed, category := createFeed(t, client)

	if feed.Title != testFeedTitle {
		t.Fatalf(`Invalid feed title, got "%v" instead of "%v"`, feed.Title, testFeedTitle)
	}

	if feed.SiteURL != testWebsiteURL {
		t.Fatalf(`Invalid site URL, got "%v" instead of "%v"`, feed.SiteURL, testWebsiteURL)
	}

	if feed.FeedURL != testFeedURL {
		t.Fatalf(`Invalid feed URL, got "%v" instead of "%v"`, feed.FeedURL, testFeedURL)
	}

	if feed.Category.ID != category.ID {
		t.Fatalf(`Invalid feed category ID, got "%v" instead of "%v"`, feed.Category.ID, category.ID)
	}

	if feed.Category.UserID != category.UserID {
		t.Fatalf(`Invalid feed category user ID, got "%v" instead of "%v"`, feed.Category.UserID, category.UserID)
	}

	if feed.Category.Title != category.Title {
		t.Fatalf(`Invalid feed category title, got "%v" instead of "%v"`, feed.Category.Title, category.Title)
	}
}

func TestGetFeedIcon(t *testing.T) {
	client := createClient(t)
	feed, _ := createFeed(t, client)
	feedIcon, err := client.FeedIcon(feed.ID)
	if err != nil {
		t.Fatal(err)
	}

	if feedIcon.ID == 0 {
		t.Fatalf(`Invalid feed icon ID, got "%v"`, feedIcon.ID)
	}

	if feedIcon.MimeType != "image/x-icon" {
		t.Fatalf(`Invalid feed icon mime type, got "%v" instead of "%v"`, feedIcon.MimeType, "image/x-icon")
	}

	if !strings.Contains(feedIcon.Data, "image/x-icon") {
		t.Fatalf(`Invalid feed icon data, got "%v"`, feedIcon.Data)
	}
}

func TestGetFeedIconNotFound(t *testing.T) {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	_, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	client = miniflux.NewClient(testBaseURL, username, testStandardPassword)
	if _, err := client.FeedIcon(42); err == nil {
		t.Fatalf(`The feed icon should be null`)
	}
}

func TestGetFeeds(t *testing.T) {
	client := createClient(t)
	feed, category := createFeed(t, client)

	feeds, err := client.Feeds()
	if err != nil {
		t.Fatal(err)
	}

	if len(feeds) != 1 {
		t.Fatalf(`Invalid number of feeds`)
	}

	if feeds[0].ID != feed.ID {
		t.Fatalf(`Invalid feed ID, got "%v" instead of "%v"`, feeds[0].ID, feed.ID)
	}

	if feeds[0].Title != testFeedTitle {
		t.Fatalf(`Invalid feed title, got "%v" instead of "%v"`, feeds[0].Title, testFeedTitle)
	}

	if feeds[0].SiteURL != testWebsiteURL {
		t.Fatalf(`Invalid site URL, got "%v" instead of "%v"`, feeds[0].SiteURL, testWebsiteURL)
	}

	if feeds[0].FeedURL != testFeedURL {
		t.Fatalf(`Invalid feed URL, got "%v" instead of "%v"`, feeds[0].FeedURL, testFeedURL)
	}

	if feeds[0].Category.ID != category.ID {
		t.Fatalf(`Invalid feed category ID, got "%v" instead of "%v"`, feeds[0].Category.ID, category.ID)
	}

	if feeds[0].Category.UserID != category.UserID {
		t.Fatalf(`Invalid feed category user ID, got "%v" instead of "%v"`, feeds[0].Category.UserID, category.UserID)
	}

	if feeds[0].Category.Title != category.Title {
		t.Fatalf(`Invalid feed category title, got "%v" instead of "%v"`, feeds[0].Category.Title, category.Title)
	}
}

func TestGetAllFeedEntries(t *testing.T) {
	client := createClient(t)
	feed, _ := createFeed(t, client)

	allResults, err := client.FeedEntries(feed.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if allResults.Total == 0 {
		t.Fatal(`Invalid number of entries`)
	}

	if allResults.Entries[0].Title == "" {
		t.Fatal(`Invalid entry title`)
	}

	filteredResults, err := client.FeedEntries(feed.ID, &miniflux.Filter{Limit: 1, Offset: 5})
	if err != nil {
		t.Fatal(err)
	}

	if allResults.Total != filteredResults.Total {
		t.Fatal(`Total should always contains the total number of items regardless of filters`)
	}

	if allResults.Entries[0].ID == filteredResults.Entries[0].ID {
		t.Fatal(`Filtered entries should be different than previous results`)
	}

	filteredResultsByEntryID, err := client.FeedEntries(feed.ID, &miniflux.Filter{BeforeEntryID: allResults.Entries[0].ID})
	if err != nil {
		t.Fatal(err)
	}

	if filteredResultsByEntryID.Entries[0].ID == allResults.Entries[0].ID {
		t.Fatal(`The first entry should filtered out`)
	}
}

func TestGetAllEntries(t *testing.T) {
	client := createClient(t)
	createFeed(t, client)

	resultWithoutSorting, err := client.Entries(nil)
	if err != nil {
		t.Fatal(err)
	}

	if resultWithoutSorting.Total == 0 {
		t.Fatal(`Invalid number of entries`)
	}

	resultWithStatusFilter, err := client.Entries(&miniflux.Filter{Status: miniflux.EntryStatusRead})
	if err != nil {
		t.Fatal(err)
	}

	if resultWithStatusFilter.Total != 0 {
		t.Fatal(`We should have 0 read entries`)
	}

	resultWithDifferentSorting, err := client.Entries(&miniflux.Filter{Order: "published_at", Direction: "desc"})
	if err != nil {
		t.Fatal(err)
	}

	if resultWithDifferentSorting.Entries[0].Title == resultWithoutSorting.Entries[0].Title {
		t.Fatalf(`The items should be sorted differently "%v" vs "%v"`, resultWithDifferentSorting.Entries[0].Title, resultWithoutSorting.Entries[0].Title)
	}

	resultWithStarredEntries, err := client.Entries(&miniflux.Filter{Starred: true})
	if err != nil {
		t.Fatal(err)
	}

	if resultWithStarredEntries.Total != 0 {
		t.Fatalf(`We are not supposed to have starred entries yet`)
	}
}

func TestSearchEntries(t *testing.T) {
	client := createClient(t)
	categories, err := client.Categories()
	if err != nil {
		t.Fatal(err)
	}

	feedID, err := client.CreateFeed("https://github.com/miniflux/miniflux/releases.atom", categories[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	if feedID == 0 {
		t.Fatalf(`Invalid feed ID, got %q`, feedID)
	}

	results, err := client.Entries(&miniflux.Filter{Search: "2.0.8"})
	if err != nil {
		t.Fatal(err)
	}

	if results.Total != 1 {
		t.Fatalf(`We should have only one entry instead of %d`, results.Total)
	}
}

func TestInvalidFilters(t *testing.T) {
	client := createClient(t)
	createFeed(t, client)

	_, err := client.Entries(&miniflux.Filter{Status: "invalid"})
	if err == nil {
		t.Fatal(`Using invalid status should raise an error`)
	}

	_, err = client.Entries(&miniflux.Filter{Direction: "invalid"})
	if err == nil {
		t.Fatal(`Using invalid direction should raise an error`)
	}

	_, err = client.Entries(&miniflux.Filter{Order: "invalid"})
	if err == nil {
		t.Fatal(`Using invalid order should raise an error`)
	}
}

func TestGetEntry(t *testing.T) {
	client := createClient(t)
	createFeed(t, client)

	result, err := client.Entries(&miniflux.Filter{Limit: 1})
	if err != nil {
		t.Fatal(err)
	}

	entry, err := client.FeedEntry(result.Entries[0].FeedID, result.Entries[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	if entry.ID != result.Entries[0].ID {
		t.Fatal("Wrong entry returned")
	}

	entry, err = client.Entry(result.Entries[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	if entry.ID != result.Entries[0].ID {
		t.Fatal("Wrong entry returned")
	}
}

func TestUpdateStatus(t *testing.T) {
	client := createClient(t)
	createFeed(t, client)

	result, err := client.Entries(&miniflux.Filter{Limit: 1})
	if err != nil {
		t.Fatal(err)
	}

	err = client.UpdateEntries([]int64{result.Entries[0].ID}, miniflux.EntryStatusRead)
	if err != nil {
		t.Fatal(err)
	}

	entry, err := client.Entry(result.Entries[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	if entry.Status != miniflux.EntryStatusRead {
		t.Fatal("The entry status should be updated")
	}

	err = client.UpdateEntries([]int64{result.Entries[0].ID}, "invalid")
	if err == nil {
		t.Fatal(`Invalid entry status should ne be accepted`)
	}
}

func TestToggleBookmark(t *testing.T) {
	client := createClient(t)
	createFeed(t, client)

	result, err := client.Entries(&miniflux.Filter{Limit: 1})
	if err != nil {
		t.Fatal(err)
	}

	if result.Entries[0].Starred {
		t.Fatal("The entry should not be starred")
	}

	err = client.ToggleBookmark(result.Entries[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	entry, err := client.Entry(result.Entries[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	if !entry.Starred {
		t.Fatal("The entry should be starred")
	}
}

func getRandomUsername() string {
	rand.Seed(time.Now().UnixNano())
	var suffix []string
	for i := 0; i < 10; i++ {
		suffix = append(suffix, strconv.Itoa(rand.Intn(1000)))
	}
	return "user" + strings.Join(suffix, "")
}

func createClient(t *testing.T) *miniflux.Client {
	username := getRandomUsername()
	client := miniflux.NewClient(testBaseURL, testAdminUsername, testAdminPassword)
	_, err := client.CreateUser(username, testStandardPassword, false)
	if err != nil {
		t.Fatal(err)
	}

	return miniflux.NewClient(testBaseURL, username, testStandardPassword)
}

func createFeed(t *testing.T, client *miniflux.Client) (*miniflux.Feed, *miniflux.Category) {
	categories, err := client.Categories()
	if err != nil {
		t.Fatal(err)
	}

	feedID, err := client.CreateFeed(testFeedURL, categories[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	if feedID == 0 {
		t.Fatalf(`Invalid feed ID, got %q`, feedID)
	}

	feed, err := client.Feed(feedID)
	if err != nil {
		t.Fatal(err)
	}

	return feed, categories[0]
}
