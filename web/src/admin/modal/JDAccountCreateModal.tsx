import { Button, Form, FormProps, Input, message, Modal } from "antd";
import { useApis } from "../../api/api";
import { IJDAccountCreate, JDAccountCreateReq } from "../../api/types";
import { useState } from "react";
import axios from "axios";

const { TextArea } = Input;

interface ModalDataType {
  isModalOpen: boolean
  onOk: Function;
  onCancel: Function;
}

type CreateFieldType = {
  accounts: string;
  type: number;
  guobai: number;
};

const JDAccountCreateModal = (params: ModalDataType) => {
  let { jdAccountCreate } = useApis()
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [formDisabled, setFormDisabled] = useState<boolean>(false);
  const [form] = Form.useForm();

  const onFinish: FormProps<CreateFieldType>['onFinish'] = async (value) => {
    setFormDisabled(true)
    setConfirmLoading(true)

    try {
      handleCreateAccounts(value)
    } catch (e) {
      if (axios.isAxiosError(e)) {
        let msg = e.response?.data?.message
        msg && message.error(msg)
      }
    } finally {
      setFormDisabled(false)
      setConfirmLoading(false)
    }
  };

  const handleCreateAccounts = async (value: CreateFieldType) => {
    let accountList: IJDAccountCreate[] = []
    if (value.accounts) {
      let pinReg = /pin=([\s\S]*?);/
      let keyReg = /wskey=([\s\S]*?);/

      let list = value.accounts.split(/[(\r\n)\r\n]+/)

      list?.forEach((line: string) => {
        let pinMatch = line.trim().match(pinReg)
        let pin = ''
        if (pinMatch && pinMatch.length > 1) {
          pin = pinMatch[1]
        }
        let keyMatch = line.trim().match(keyReg)
        let wsKey = ''
        if (keyMatch && keyMatch.length > 1) {
          wsKey = keyMatch[1]
        }

        if (pin && wsKey) {
          accountList.push({
            account: pin,
            wsKey: wsKey,
          })
        }
      })
    }

    let data: JDAccountCreateReq = {
      accountList: accountList,
      remark: ""
    }

    await jdAccountCreate(data)

    params?.onOk?.();
    message.success('导入成功')
  }

  return (
    <>
      <Modal
        title="导入京东账号"
        okText="确认"
        cancelText="取消"
        onCancel={() => {
          form.resetFields();
          params?.onCancel?.();
        }}
        confirmLoading={confirmLoading}
        open={params.isModalOpen}
        // 添加宽度限制
        width={500}
        style={{ maxWidth: '90%' }}
        footer={[
          <Button key="cancel" onClick={() => {
            form.resetFields();
            params?.onCancel?.();
          }}>取消</Button>,
          <Button
            key="submit"
            type="primary"
            loading={confirmLoading}
            onClick={() => {
              form.submit();
            }}
          >
            确认
          </Button>,
        ]}
      >
        <Form
          form={form}
          layout="vertical"
          disabled={formDisabled}
          onFinish={onFinish}
        >
          {/* <Form.Item<CreateFieldType>
            name="type"
            label="类型"
          >
            <Select options={[
              { value: '1', label: '游戏实物' },
            ]} style={{ width: 200 }}>
            </Select>
          </Form.Item>

          <Form.Item<CreateFieldType>
            name="guobai"
            label="是否国白或复登（只下一单）"
          >
            <Select options={[
              { value: '1', label: '是' },
              { value: '2', label: '否' },
            ]} style={{ width: 200 }}>
            </Select>
          </Form.Item> */}

          <Form.Item<CreateFieldType>
            name="accounts"
            label="账号"
          >
            <TextArea rows={8} />
          </Form.Item>
        </Form >
      </Modal>
    </>
  );
};

export default JDAccountCreateModal;