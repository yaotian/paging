package paging

import (
	"errors"
	"github.com/yaotian/logs"
	"strconv"
)

const (
	lineSize          = 10
	show_HowMany_Page = 4
)

type Paging interface {
	SetCurrentPage(uint)
	CurrentPage() uint
	LineSize() uint
	SetTotalPage()
	SetPageScope(int)
	TotalPage() uint
	TotalCount() uint
	PageScope() []int
}

type paging struct {
	currentPage uint  // 当前页
	lineSize    uint  // 每页几行
	totalPage   uint  // 总页数
	totalCount  uint  // 总行数
	pageScope   []int //下面页面的范围
}

func New(lineSize, totalCount uint) Paging {
	return &paging{
		lineSize:   lineSize,
		totalCount: totalCount,
	}
}

func (p *paging) SetCurrentPage(currentPage uint) {
	p.currentPage = currentPage
}

func (p paging) CurrentPage() uint {
	return p.currentPage
}

func (p paging) LineSize() uint {
	return p.lineSize
}

func (p *paging) SetTotalPage() {
	totalPage := p.TotalCount() / p.LineSize()
	if p.TotalCount()%p.LineSize() == 0 {
		p.totalPage = totalPage
	} else {
		p.totalPage = totalPage + 1
	}
}

func (p paging) TotalPage() uint {
	return p.totalPage
}

func (p paging) TotalCount() uint {
	return p.totalCount
}

func (p paging) PageScope() []int {
	return p.pageScope
}

func (p *paging) SetPageScope(current_page int) {
	result := []int{}
	int_totalPage := int(p.totalPage)
	b := 0
	if show_HowMany_Page >= int_totalPage { // item is very slow and is less than show_HowMany_Page
		b = int_totalPage
		current_page = 1
	} else if current_page+show_HowMany_Page > int_totalPage { //选择的比较靠后的page
		b = int_totalPage
		current_page = int_totalPage - show_HowMany_Page
	} else {
		b = current_page + show_HowMany_Page
	}
	for ; current_page <= b; current_page++ {
		result = append(result, current_page)
	}
	p.pageScope = result
}

func Pagination(current_page_str string, totalCount uint) (Paging, error) {

	// New Paging
	page := New(lineSize, totalCount)

	// set total page
	page.SetTotalPage()

	currentPage := uint(0)
	if len(current_page_str) == 0 {
		currentPage = uint(1)
	} else {
		newCurrentPage, err := strconv.ParseUint(current_page_str, 10, 32)
		if err != nil {
			logs.Logger.Debug("Atoi: ", err.Error())
			return nil, err
		}
		if newCurrentPage < 1 {
			currentPage = uint(1)
		} else if uint(newCurrentPage) >= page.TotalCount() {
			currentPage = page.TotalCount()
		}
		currentPage = uint(newCurrentPage)
	}
	page.SetPageScope(int(currentPage))
	page.SetCurrentPage(currentPage)

	return page, nil
}

func Make_paging(current_page_str string, data_for_process []interface{}) (interface{}, error) {

	len_of_coming_data := len(data_for_process)

	page, err := Pagination(current_page_str, uint(len_of_coming_data))
	if err != nil {
		logs.Logger.Error("paging: ", err.Error())
		return nil, errors.New("error")
	}

	from_current_page_str, err := strconv.ParseUint(current_page_str, 10, 32)
	if err != nil {
		logs.Logger.Error("Atoi: ", err.Error())
	}

	if int(from_current_page_str) > int(page.TotalPage()) {
		return nil, errors.New("current page is bigger than total page. wrong request")
	}

	logs.Logger.Debug(current_page_str)
	logs.Logger.Debug(page.CurrentPage())
	logs.Logger.Debug(page.TotalPage())

	firstResult := (page.CurrentPage() - uint(1)) * page.LineSize()
	maxResult := page.LineSize()

	result_for_show := data_for_process

	if firstResult+maxResult < uint(len_of_coming_data) {
		result_for_show = (data_for_process)[firstResult : firstResult+maxResult]
	} else {
		logs.Logger.Debug(firstResult)
		logs.Logger.Debug(len_of_coming_data)
		result_for_show = (data_for_process)[firstResult:len_of_coming_data]
	}

	prev := page.CurrentPage() - 1
	if prev <= 1 {
		prev = 1
	}

	next := page.CurrentPage() + 1
	if next >= page.TotalPage() {
		next = page.TotalPage()
	}

	var data = struct {
		Results     []interface{}
		CurrentPage string
		TotalPage   string
		Prev        string
		Next        string
		PageScope   []int
	}{
		Results:     result_for_show,
		CurrentPage: strconv.Itoa(int(page.CurrentPage())),
		TotalPage:   strconv.Itoa(int(page.TotalPage())),
		Prev:        strconv.Itoa(int(prev)),
		Next:        strconv.Itoa(int(next)),
		PageScope:   page.PageScope(),
	}
	//for gc
	result_for_show = nil

	return data, nil

	//return renderTemplate("index", data)
}
