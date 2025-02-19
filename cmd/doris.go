// Package cmd @Author lanpang
// @Date 2025/1/14 下午3:34:00
// @Desc
package cmd

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"regexp"
	"strings"
)

type Doris struct {
	ProjectName        string
	InspectionDate     string
	JobFailureCount    int
	StaffCount         int
	UseAnalyseCount    int
	CustomerGroupCount int
	OnlineBENodes      int
	TotalBENodes       int
}

// DorisHandler Doris 数据处理器
type DorisHandler struct {
	reProjectName        *regexp.Regexp
	reInspectionDate     *regexp.Regexp
	reTotalBENodes       *regexp.Regexp
	reOnlineBENodes      *regexp.Regexp
	reJobFailure         *regexp.Regexp
	reStaffCount         *regexp.Regexp
	reUseAnalyseCount    *regexp.Regexp
	reCustomerGroupCount *regexp.Regexp
}

// NewDorisHandler 初始化处理器
func NewDorisHandler() *DorisHandler {
	return &DorisHandler{
		reProjectName:        regexp.MustCompile(`\*\*项目名称：.*>(.*?)</font>`),
		reInspectionDate:     regexp.MustCompile(`\*\*巡检时间：.*>(.*?)</font>`),
		reTotalBENodes:       regexp.MustCompile(`\*\*BE节点总数：.*>(\d+)</font>`),
		reOnlineBENodes:      regexp.MustCompile(`\*\*在线节点数：.*>(\d+)</font>`),
		reJobFailure:         regexp.MustCompile(`\*\*Job失败数：.*>(\d+)</font>`),
		reStaffCount:         regexp.MustCompile(`\*\*员工统计表：.*>(\d+)</font>`),
		reUseAnalyseCount:    regexp.MustCompile(`\*\*使用分析表：.*>(\d+)</font>`),
		reCustomerGroupCount: regexp.MustCompile(`\*\*客户群统计表：.*>(\d+)</font>`),
	}
}

func (h *DorisHandler) Handle(content string) *Doris {
	// 提取字段
	projectName := extractMatch(h.reProjectName, content)
	inspectionDate := extractMatch(h.reInspectionDate, content)
	totalBENodes := parseInt(extractMatch(h.reTotalBENodes, content))
	onlineBENodes := parseInt(extractMatch(h.reOnlineBENodes, content))
	jobFailure := parseInt(extractMatch(h.reJobFailure, content))
	staffCount := parseInt(extractMatch(h.reStaffCount, content))
	useAnalyse := parseInt(extractMatch(h.reUseAnalyseCount, content))
	customerGroupCount := parseInt(extractMatch(h.reCustomerGroupCount, content))
	// 打印staffCount、useAnalyse、customerGroupCount
	//log.Println("staffCount:", staffCount)
	//log.Println("useAnalyse:", useAnalyse)
	//log.Println("customerGroupCount:", customerGroupCount)

	// 构造 DorisData 结构体
	return &Doris{
		ProjectName:        projectName,
		InspectionDate:     inspectionDate,
		TotalBENodes:       totalBENodes,
		OnlineBENodes:      onlineBENodes,
		JobFailureCount:    jobFailure,
		StaffCount:         staffCount,
		UseAnalyseCount:    useAnalyse,
		CustomerGroupCount: customerGroupCount,
	}
}

func insertDorisLog(d *Doris, pgClient *pgx.Conn) {
	// 查询项目ID
	projectId := selectProjectId(pgClient, d.ProjectName)
	// 插入数据
	insertQuery := `
		INSERT INTO public.qw_project_doris_log 
		    (project_name, project_id, job_failure_count, staff_count, use_analyse_count, customer_group_count, online_be_nodes, total_be_nodes, inspection_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (project_name, inspection_date) DO UPDATE
		SET
		    project_id = EXCLUDED.project_id,
		    job_failure_count = EXCLUDED.job_failure_count,
		    staff_count = EXCLUDED.staff_count,
		    use_analyse_count = EXCLUDED.use_analyse_count,
		    customer_group_count = EXCLUDED.customer_group_count,
		    online_be_nodes = EXCLUDED.online_be_nodes,
		    total_be_nodes = EXCLUDED.total_be_nodes
		RETURNING id;`
	_, err := pgClient.Exec(context.Background(), insertQuery,
		d.ProjectName, projectId, d.JobFailureCount, d.StaffCount, d.UseAnalyseCount, d.CustomerGroupCount, d.OnlineBENodes, d.TotalBENodes, d.InspectionDate)
	if err != nil {
		log.Println(d.ProjectName, "插入数据失败: ", err)
	} else {
		//log.Println(d.ProjectName, "数据插入成功!")
	}
}

func checkDorisData(d *Doris, pgClient *pgx.Conn) {
	isAlert := false
	if d.JobFailureCount > 0 {
		log.Println(d.ProjectName, "Job失败数异常!")
		isAlert = true
	}
	if d.StaffCount == 0 && d.UseAnalyseCount == 0 && d.CustomerGroupCount == 0 {
		log.Println(d.ProjectName, "员工、使用分析、客户群统计表异常!")
		isAlert = true
	}
	if d.OnlineBENodes < d.TotalBENodes {
		log.Println(d.ProjectName, "BE节点数异常!")
		isAlert = true
	}
	if isAlert {
		// 查询项目运维
		operationer := selectProjectOperationer(pgClient, selectProjectId(pgClient, d.ProjectName))
		toList := getToList(operationer)
		var builder strings.Builder
		// 构建邮件内容
		builder.WriteString(dorisMailHead)
		// 构建表格行
		dorisData := fmt.Sprintf("<td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td></tr>", d.ProjectName, d.InspectionDate, d.JobFailureCount, d.StaffCount, d.UseAnalyseCount, d.CustomerGroupCount, d.TotalBENodes, d.OnlineBENodes)
		builder.WriteString(dorisData)
		builder.WriteString("</table><br>请查收！</br><br>Send By notify@wshoto.com </br>（自动发送请勿回复）")
		mailAlert("Doris", builder.String(), toList)
	}
}
