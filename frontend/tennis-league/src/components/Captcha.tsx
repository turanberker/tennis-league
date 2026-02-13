import React, { useState, useEffect, useRef } from 'react';
import { InputText } from 'primereact/inputtext';

interface CaptchaProps {
  value: string;
  onChange: (val: string) => void;
}

const Captcha: React.FC<CaptchaProps> = ({ value, onChange }) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [captchaText, setCaptchaText] = useState<string>('');

  const generateCaptcha = () => {
    const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789';
    let text = '';
    for (let i = 0; i < 5; i++) {
      text += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    setCaptchaText(text);
  };

  useEffect(() => {
    generateCaptcha();
  }, []);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.fillStyle = '#f2f2f2';
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    ctx.font = '24px Arial';
    ctx.fillStyle = '#000';
    ctx.textBaseline = 'middle';
    ctx.textAlign = 'center';
    ctx.fillText(captchaText, canvas.width / 2, canvas.height / 2);
  }, [captchaText]);

  return (
    <div className="flex flex-column gap-2">
      <canvas
        ref={canvasRef}
        width={150}
        height={50}
        style={{ border: '1px solid #ccc', borderRadius: '4px' }}
        onClick={generateCaptcha} // tıklayınca yeniler
      />
      <InputText
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder="Captcha kodunu girin"
      />
      <small>
        Captcha kodunu yukarıdaki resimden girin. (Click to refresh)
      </small>
    </div>
  );
};

export default Captcha;
