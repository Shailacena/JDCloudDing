import { useCallback, useEffect, useState } from "react";
import { Form, Table, Input, Button, Card, Divider, message, Space, Flex, ConfigProvider, Modal, DatePicker } from "antd";
import type { FormProps, TableProps } from "antd";
import { useApis } from "../api/api";
import { IJDAccount, JDAccountDeleteReq, JDAccountResetReq, JDAccountSearchParams, ListJDAccountReq } from "../api/types";
import axios from "axios";
import { getDataFormat } from "../utils/Tool";
import JDAccountCreateModal from "./modal/JDAccountCreateModal";
import { SearchOutlined } from "@ant-design/icons";
import { convertJDAccountStatus, JDAccountStatus } from "../utils/biz";
import { PAGE_DEFAULT_INDEX, PAGE_SIZE } from "../components/types";
import { Dayjs } from "dayjs";

interface DataType extends IJDAccount {
  key: number;
}

export interface IJDAccountSearchCondition extends JDAccountSearchParams {
  dayRange: Dayjs[]
}

function JDAccount() {
  const columns: TableProps<DataType>["columns"] = [
    {
      title: "ID",
      dataIndex: "id",
      key: "id",
    },
    {
      title: "账号",
      dataIndex: "account",
      key: "account",
      width: "300px",
      ellipsis: true,
    },
    {
      title: "状态",
      key: "status",
      dataIndex: "status",
      render: (_, d) => (
        <span
          style={{
            color: d.status === JDAccountStatus.Normal ? "#52c41a" : "#f5222d",
          }}
        >
          {convertJDAccountStatus(d.status)}
        </span>
      ),
    },
    {
      title: "创建时间",
      key: "createAt",
      dataIndex: "createAt",
      render: (ts: number) => {
        return getDataFormat(new Date(ts * 1000));
      },
    },
    {
      title: "备注",
      key: "remark",
      dataIndex: "remark",
    },
    {
      title: "操作",
      key: "action",
      fixed: "right",
      render: (_, d) => {
        let enable = d.status === JDAccountStatus.Normal;
        return (
          <Space size="middle">
            <Button
              type="primary"
              size="small"
              danger={enable}
              onClick={() => {
                enableAcount(
                  d.id,
                  enable ? JDAccountStatus.Invalid : JDAccountStatus.Normal
                );
              }}
            >
              {enable ? "禁用" : "启用"}
            </Button>
          </Space>
        );
      },
    },
  ];

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [list, setList] = useState<DataType[]>([]);
  let { listJDAccount, jdAccountEnable, jdAccountDelete, jdAccountReset } =
    useApis();
  const [total, setTotal] = useState(0);
  const [reqParams, setReqParams] = useState<ListJDAccountReq>({
    currentPage: PAGE_DEFAULT_INDEX,
    pageSize: PAGE_SIZE.TEN,
  });
  const [listLoadingStates, setListLoadingStates] = useState(false);

  const showModal = () => {
    setIsModalOpen(true);
  };

  const handleOk = () => {
    setIsModalOpen(false);
    fetchJDAccountList();
  };

  const handleCancel = () => {
    setIsModalOpen(false);
  };

  const enableAcount = async (id: number, status: number) => {
    try {
      await jdAccountEnable({ id, status });
      fetchJDAccountList();
      message.success(
        status === JDAccountStatus.Normal ? "启用成功" : "禁用成功"
      );
    } catch (e) {
      if (axios.isAxiosError(e)) {
        let msg = e.response?.data?.message;
        msg && message.error(msg);
      }
    }
  };

  const fetchJDAccountList = useCallback(async () => {
    try {
      setListLoadingStates(true);

      const { data } = await listJDAccount(reqParams);
      let d: DataType[] = data?.list?.map((item, index) => {
        let newItem: DataType = {
          key: index,
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
    fetchJDAccountList();
  }, [reqParams]);

  const onSearch: FormProps<IJDAccountSearchCondition>["onFinish"] = async (
    value
  ) => {
    let { dayRange} = value
    setReqParams({ ...value, currentPage: 1, pageSize: reqParams.pageSize, startAt: dayRange?.[0]?.format('YYYY-MM-DD'),
      endAt: dayRange?.[1]?.format('YYYY-MM-DD') });
  };

  const handleTableChange = (current: number, pageSize: number) => {
    setReqParams({ ...reqParams, currentPage: current, pageSize });
  };

  const checkRemove = (isAll: boolean) => {
    const {
      currentPage: currentPage,
      pageSize: pageSize,
      ...params
    } = reqParams;

    let msg = "确定删除全部ck？";
    if (!isAll) {
      msg = "确定删除指定搜索条件ck？";
      if (!params.id && !params.account) {
        message.error("请指定搜索条件");
        return;
      }
    }

    Modal.confirm({
      content: msg,
      onOk: () => {
        remove(isAll);
      },
    });
  };

  const remove = async (isAll: boolean) => {
    const {
      currentPage: currentPage,
      pageSize: pageSize,
      ...params
    } = reqParams;

    let p: JDAccountDeleteReq = { isAll };
    if (!isAll) {
      p = { ...p, ...params, id: params?.id && +params?.id };
    }

    try {
      await jdAccountDelete(p);
      message.success("删除成功");
      fetchJDAccountList();
    } catch (e) {
      console.error(e);
    }
  };

  const checkResetStatus = async (idx: number) => {
    let msg = "确定重置所有失败ck？";
    let status: number[] = [];

    switch (idx) {
      case 1:
        msg = "确定重置所有失败ck？";
        status = [3, 4, 5, 6];
        break;
      case 2:
        msg = "确定重置所有过期ck？";
        status = [2];
        break;
      case 3:
        const {
          currentPage: currentPage,
          pageSize: pageSize,
          ...params
        } = reqParams;

        if (!params.id && !params.account) {
          message.error("请指定搜索条件");
          return;
        }

        msg = "确定重置指定搜索条件的ck？";
        status = [2, 3, 4, 5, 6];
        break;
    }

    Modal.confirm({
      content: msg,
      onOk: () => {
        reset(idx, status);
      },
    });
  };

  const reset = async (idx: number, status: Array<number>) => {
    const {
      currentPage: currentPage,
      pageSize: pageSize,
      ...params
    } = reqParams;

    let p: JDAccountResetReq = { ...params, status: status };
    if (idx === 3 && reqParams.id) {
      p.id = +reqParams.id;
    }

    try {
      await jdAccountReset(p);
      message.success("重置成功");
      fetchJDAccountList();
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <>
      <Card>
        <Form
          className="inline_search_form"
          name="inline_search_form"
          layout="inline"
          onFinish={onSearch}
        >
          <Form.Item<ListJDAccountReq> name="id">
            <Input allowClear placeholder="ID" />
          </Form.Item>

          <Form.Item<ListJDAccountReq> name="account">
            <Input allowClear placeholder="账号" />
          </Form.Item>

          <Form.Item<IJDAccountSearchCondition> name="dayRange">
            <DatePicker.RangePicker style={{ width: 250 }} />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              icon={<SearchOutlined />}
              htmlType="submit"
            ></Button>
          </Form.Item>
        </Form>

        <Divider />
        <ConfigProvider>
          <Flex gap="small" wrap>
            <Button type="primary" onClick={showModal}>
              批量导入京东账号
            </Button>
            <Button
              variant="solid"
              color="green"
              onClick={() => {
                checkResetStatus(1);
              }}
            >
              重置失败ck
            </Button>
            <Button
              variant="solid"
              color="orange"
              onClick={() => {
                checkResetStatus(2);
              }}
            >
              重置过期ck
            </Button>
            <Button
              variant="solid"
              color="magenta"
              onClick={() => {
                checkResetStatus(3);
              }}
            >
              重置指定搜索条件ck
            </Button>
            <Button
              type="primary"
              danger
              onClick={() => {
                checkRemove(false);
              }}
            >
              删除指定搜索条件ck
            </Button>
            <Button
              type="primary"
              danger
              onClick={() => {
                checkRemove(true);
              }}
            >
              全部删除
            </Button>
          </Flex>
        </ConfigProvider>
        <Divider />
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

        {isModalOpen && (
          <JDAccountCreateModal
            isModalOpen={isModalOpen}
            onOk={handleOk}
            onCancel={handleCancel}
          />
        )}
      </Card>
    </>
  );
}

export default JDAccount;
