import { useCallback, useEffect, useState } from "react";
import { Card, Table } from "antd";
import type { TableProps } from "antd";
import { useApis } from "../api/api";
import { IOperationLog, ListOperationLogReq } from "../api/types";
import { convertTimestamp } from "../utils/biz";
import { PAGE_DEFAULT_INDEX, PAGE_SIZE } from "../components/types";

interface DataType extends IOperationLog {
  key: string;
}

const columns: TableProps<DataType>["columns"] = [
  {
    title: "操作者ID",
    dataIndex: "operator",
    key: "operator",
    align: "center",
  },
  {
    title: "操作者名称",
    key: "operatorName",
    dataIndex: "operatorName",
    align: "center",
  },
  {
    title: "内容",
    key: "operationStr",
    dataIndex: "operationStr",
    align: "center",
  },
  {
    title: "ip",
    key: "ip",
    dataIndex: "ip",
    align: "center",
  },
  {
    title: "时间",
    key: "createdAt",
    dataIndex: "createdAt",
    align: "center",
    render: (_, d) => {
      return convertTimestamp(d.createAt);
    },
  },
];

function OperationLog() {
  const [list, setList] = useState<DataType[]>([]);
  let { listOperationLog } = useApis();
  const [reqParams, setReqParams] = useState<ListOperationLogReq>({
    currentPage: PAGE_DEFAULT_INDEX,
    pageSize: PAGE_SIZE.TEN,
  });
  const [listLoadingStates, setListLoadingStates] = useState(false);
  const [total, setTotal] = useState(0);

  const fetchListOperationLog = useCallback(async () => {
    try {
      setListLoadingStates(true);

      const { data } = await listOperationLog(reqParams);
      let d: DataType[] = data?.list?.map((item, index) => {
        let newItem: DataType = {
          key: index.toString(),
          ...item,
        };
        return newItem;
      });

      setList(d);
      setTotal(data.total);
    } catch (e) {
      console.error(e);
    } finally {
      setListLoadingStates(false);
    }
  }, [reqParams]);

  useEffect(() => {
    fetchListOperationLog();
  }, [reqParams]);

  const handleTableChange = (current: number, pageSize: number) => {
    setReqParams({ ...reqParams, currentPage: current, pageSize });
  };

  return (
    <>
      <Card>
        <Table<DataType>
          bordered
          size="middle"
          pagination={{
            current: reqParams.currentPage,
            pageSize: reqParams.pageSize,
            total: total,
            onChange: handleTableChange,
          }}
          columns={columns}
          dataSource={list}
          scroll={{ x: "max-content" }}
          loading={listLoadingStates}
        />
      </Card>
    </>
  );
}

export default OperationLog;
