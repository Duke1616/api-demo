package impl

import (
	"context"
	"fmt"
	"github.com/Duke1616/api-demo/api/host"
	"github.com/infraboard/mcube/sqlbuilder"
)

var _ host.Service = &impl{}

// CreateHost 把Host对象保存到数据内
func (i *impl) CreateHost(ctx context.Context, ins *host.Host) error {
	var (
		err error
	)

	// 把数据入库到 resource表和host表
	// 一次需要往2个表录入数据, 我们需要2个操作 要么都成功，要么都失败, 事务的逻辑
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start tx error, %s", err)
	}

	// 通过Defer处理事务提交方式
	// 1. 无报错，则Commit 事务
	// 2. 有报错, 则Rollback 事务
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				i.log.Error("rollback error, %s", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				i.log.Error("commit error, %s", err)
			}
		}
	}()

	// 插入Resource数据
	resStmt, err := tx.Prepare(insertResourceSQL)
	if err != nil {
		return err
	}
	defer resStmt.Close()

	_, err = resStmt.Exec(
		ins.Id, ins.Vendor, ins.Region, ins.Zone, ins.CreateAt, ins.ExpireAt, ins.Category, ins.Type, ins.InstanceId,
		ins.Name, ins.Description, ins.Status, ins.UpdateAt, ins.SyncAt, ins.SyncAccount, ins.PublicIP,
		ins.PrivateIP, ins.PayType, ins.ResourceHash, ins.DescribeHash,
	)
	if err != nil {
		return err
	}

	// 插入Describe 数据
	dstmt, err := tx.Prepare(insertDescribeSQL)
	if err != nil {
		return err
	}
	defer dstmt.Close()

	_, err = dstmt.Exec(
		ins.Id, ins.CPU, ins.Memory, ins.GPUAmount, ins.GPUSpec,
		ins.OSType, ins.OSName, ins.SerialNumber,
	)
	if err != nil {
		return err
	}

	return nil
}

func (i *impl) QueryHost(ctx context.Context, req *host.QueryHostRequest) (*host.Set, error) {
	query := sqlbuilder.NewQuery(queryHostSQL).Order("create_at").Desc().Limit(int64(req.Offset()), uint(req.PageNumber))

	sqlStr, args := query.BuildQuery()
	i.log.Debugf("sql: %s, args: %v", sqlStr, args)

	stmt, err := i.db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	set := host.NewSet()
	for rows.Next() {
		ins := host.NewDefaultHost()
		if err = rows.Scan(
			&ins.Id, &ins.Vendor, &ins.Region, &ins.Zone, &ins.CreateAt, &ins.ExpireAt,
			&ins.Category, &ins.Type, &ins.InstanceId, &ins.Name,
			&ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt, &ins.SyncAccount,
			&ins.PublicIP, &ins.PrivateIP, &ins.PayType, &ins.ResourceHash, &ins.DescribeHash,
			&ins.Id, &ins.CPU,
			&ins.Memory, &ins.GPUAmount, &ins.GPUSpec, &ins.OSType, &ins.OSName,
			&ins.SerialNumber, &ins.ImageID, &ins.InternetMaxBandwidthOut, &ins.InternetMaxBandwidthIn,
			&ins.KeyPairName, &ins.SecurityGroups,
		); err != nil {
			return nil, err
		}
		set.Add(ins)
	}

	// Count总数查询
	countStr, countArgs := query.BuildCount()
	countStmt, err := i.db.Prepare(countStr)
	if err != nil {
		return nil, err
	}

	defer countStmt.Close()

	if err = countStmt.QueryRow(countArgs...).Scan(&set.Total); err != nil {
		return nil, err
	}
	
	return set, nil
}

func (i *impl) DescribeHost(ctx context.Context, req *host.DescribeHostRequest) (*host.Host, error) {
	//TODO implement me
	panic("implement me")
}

func (i *impl) UpdateHost(ctx context.Context, req *host.UpdateHostRequest) (*host.Host, error) {
	//TODO implement
	panic("implement me")
}

func (i *impl) DeleteHost(ctx context.Context, req *host.DeleteHostRequest) (*host.Host, error) {
	//TODO implement me
	panic("implement me")
}
