import React, { useState } from 'react';
import { Button, Modal, Space } from 'antd';
import { createStyles } from 'antd-style';

const useStyle = createStyles(() => ({
  'my-modal-body': {
    // background: token.blue1,
    // padding: token.paddingSM,
  },
  'my-modal-mask': {
    // boxShadow: `inset 0 0 15px #fff`,
  },
  'my-modal-header': {
    // borderBottom: `1px dotted ${token.colorPrimary}`,
  },
  'my-modal-footer': {
    // color: token.colorPrimary,
  },
  'my-modal-content': {
    // border: '1px solid #333',
  },
}));

const OrderDetail: React.FC = () => {
  const [isModalOpen, setIsModalOpen] = useState([false, false]);
  const { styles } = useStyle();
  // const token = useTheme();

  const toggleModal = (idx: number, target: boolean) => {
    setIsModalOpen((p) => {
      p[idx] = target;
      return [...p];
    });
  };

  const classNames = {
    body: styles['my-modal-body'],
    mask: styles['my-modal-mask'],
    header: styles['my-modal-header'],
    footer: styles['my-modal-footer'],
    content: styles['my-modal-content'],
  };

  const modalStyles = {
    header: {
      // borderLeft: `5px solid ${token.colorPrimary}`,
      // borderRadius: 0,
      // paddingInlineStart: 5,
    },
    body: {
      // boxShadow: 'inset 0 0 5px #999',
      // borderRadius: 5,
    },
    mask: {
      // backdropFilter: 'blur(10px)',
    },
    footer: {
      // borderTop: '1px solid #333',
    },
    content: {
      // boxShadow: '0 0 30px #999',
    },
  };

  return (
    <>
      <Space>
        <Button type="primary" onClick={() => toggleModal(0, true)}>
          Open Modal
        </Button>
      </Space>
      <Modal
        title="Detail"
        open={isModalOpen[0]}
        onOk={() => toggleModal(0, false)}
        onCancel={() => toggleModal(0, false)}
        classNames={classNames}
        styles={modalStyles}
        // centered //居中
        style={{ top: 50 }}
        destroyOnClose
      >
        <div style={{ display: 'flex', minHeight: 200, margin: '0 auto', justifyContent: 'center', alignItems: 'center' }}>
          <p>Partner ID: 123</p>
          <p>Partner name: 123</p>
          <p>Partner ID: 123</p>
          <p>Partner ID: 123</p>
        </div>
      </Modal>
    </>
  );
};

export default OrderDetail;