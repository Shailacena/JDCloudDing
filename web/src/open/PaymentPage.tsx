import { useEffect, useState } from 'react';
import { Button, Card, Col, Row, Space, Typography } from 'antd';
import { WechatOutlined } from '@ant-design/icons';

const { Text } = Typography;

const PaymentPage = () => {
  const [timeLeft, setTimeLeft] = useState(60); // 60秒倒计时 (单位: 秒)

  useEffect(() => {
    setTimeout(() => {
      if (timeLeft <= 0) {
        alert('已超时');
      } else {
        setTimeLeft(timeLeft - 1);
      }
    }, 1000);
  });

  const onPaymentClick = () => {
    // 这里可以处理点击支付时的逻辑
    alert('支付已点击');
  };

  return (
    <div style={{ height: '100vh', margin: '0 auto', justifyContent: 'center', alignItems: 'center' }}>
      <div style={{ padding: '20px' }}>
        <Card bordered={false}>
          <Space style={{ justifyContent: 'center', width: '100%' }}>
            <WechatOutlined style={{ fontSize: 40, color: '#1AAD19' }} />
            <Text style={{ fontSize: '30px', textAlign: 'center', paddingLeft: '10px' }}>微信支付</Text>
          </Space>

        </Card>
      </div>

      <div style={{ padding: '20px' }}>
        <Card bordered={false}>
          <Row gutter={16}>
            <Col span={6} style={{ textAlign: 'left' }}>商品名称:</Col>
            <Col span={18} style={{ textAlign: 'right' }}>购买充值</Col>
          </Row>
          <Row gutter={16}>
            <Col span={6} style={{ textAlign: 'left' }}>商品价格:</Col>
            <Col span={18} style={{ textAlign: 'right' }}>500元</Col>
          </Row>
          <Row gutter={16}>
            <Col span={6} style={{ textAlign: 'left' }}>单号:</Col>
            <Col span={18} style={{ textAlign: 'right' }}>P1864667052451065858</Col>
          </Row>
        </Card>
      </div>
      
      <div>
      <Text style={{ display: 'block', color: 'black', textAlign: 'center', marginTop: 10 }}>
        请在<span style={{ color: 'red'}}>{Math.floor(timeLeft / 60)}</span>分钟<span style={{ color: 'red'}}>{timeLeft%60}</span>秒内完成支付
      </Text>
      </div>

      <div style={{ padding: 20 }}>
        <Button
          type="primary"
          block
          size="large"
          style={{ marginTop: 20 }}
          onClick={onPaymentClick}
        >
          点击支付
        </Button>
      </div>

      <Text style={{ display: 'block', color: 'black', textAlign: 'center', marginTop: 10 }}>
        提示：
      </Text>
      <div style={{ paddingLeft: 40, paddingRight: 40 }}>
        <Text style={{ display: 'block', color: 'red', textAlign: 'center', marginTop: 10 }}>
          支付成功后跳转页面提示支付失败属于正常，请不要重复支付，等待充值。
        </Text>
        <Text style={{ display: 'block', color: 'red', textAlign: 'center' }}>
          支付5分钟内未到账，请及时联系客服。过后未查不退补。
        </Text>
      </div>
    </div >
  );
};

export default PaymentPage;