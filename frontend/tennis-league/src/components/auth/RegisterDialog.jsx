import React from 'react';
import { Dialog } from 'primereact/dialog';
import { InputText } from 'primereact/inputtext';
import { Password } from 'primereact/password';
import { Button } from 'primereact/button';
import { classNames } from 'primereact/utils';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import Captcha from '../Captcha';

export default function RegisterDialog({ visible, onHide, onRegister }) {
  const validationSchema = Yup.object().shape({
    email: Yup.string()
      .email('Geçerli bir email giriniz')
      .required('Email zorunludur'),

    name: Yup.string().required('Ad zorunludur'),

    surname: Yup.string().required('Soyad zorunludur'),

    password: Yup.string().required('Şifre zorunludur'),

    passwordRepeat: Yup.string()
      .oneOf([Yup.ref('password'), null], 'Şifreler aynı olmalıdır')
      .required('Şifre tekrar zorunludur'),

    captchaInput: Yup.string().required('Captcha zorunludur'),
  });

  const formik = useFormik({
    initialValues: {
      email: '',
      name: '',
      surname: '',
      password: '',
      passwordRepeat: '',
      captchaInput: '',
    },
    validationSchema,
    onSubmit: (values) => {
      onRegister(values);
    },
  });

  const isFormFieldValid = (name) =>
    !!(formik.touched[name] && formik.errors[name]);

  const getFormErrorMessage = (name) =>
    isFormFieldValid(name) && (
      <small className="p-error">{formik.errors[name]}</small>
    );

  const footer = (
    <Button
      label="Kayıt Ol"
      icon="pi pi-user-plus"
      onClick={formik.handleSubmit}
    />
  );

  return (
    <Dialog
      header="Kayıt Ol"
      visible={visible}
      style={{ width: '420px' }}
      modal
      onHide={onHide}
      footer={footer}
    >
      <form onSubmit={formik.handleSubmit} className="flex flex-column gap-3">
        <span className="p-float-label">
          <InputText
            id="email"
            name="email"
            value={formik.values.email}
            onChange={formik.handleChange}
            onBlur={formik.handleBlur}
            className={classNames({ 'p-invalid': isFormFieldValid('email') })}
          />
          <label htmlFor="email">Email</label>
        </span>
        {getFormErrorMessage('email')}

        <span className="p-float-label">
          <InputText
            id="name"
            name="name"
            value={formik.values.name}
            onChange={formik.handleChange}
            onBlur={formik.handleBlur}
            className={classNames({ 'p-invalid': isFormFieldValid('name') })}
          />
          <label htmlFor="name">Ad</label>
        </span>
        {getFormErrorMessage('name')}

        <span className="p-float-label">
          <InputText
            id="surname"
            name="surname"
            value={formik.values.surname}
            onChange={formik.handleChange}
            onBlur={formik.handleBlur}
            className={classNames({ 'p-invalid': isFormFieldValid('surname') })}
          />
          <label htmlFor="surname">Soyad</label>
        </span>
        {getFormErrorMessage('surname')}

        <span className="p-float-label">
          <Password
            id="password"
            name="password"
            value={formik.values.password}
            onChange={formik.handleChange}
            onBlur={formik.handleBlur}
            toggleMask
            feedback={false}
            className={classNames({
              'p-invalid': isFormFieldValid('password'),
            })}
          />
          <label htmlFor="password">Şifre</label>
        </span>
        {getFormErrorMessage('password')}

        <span className="p-float-label">
          <Password
            id="passwordRepeat"
            name="passwordRepeat"
            value={formik.values.passwordRepeat}
            onChange={formik.handleChange}
            onBlur={formik.handleBlur}
            toggleMask
            feedback={false}
            className={classNames({
              'p-invalid': isFormFieldValid('passwordRepeat'),
            })}
          />
          <label htmlFor="passwordRepeat">Şifre Tekrar</label>
        </span>
        {getFormErrorMessage('passwordRepeat')}

        <Captcha
          value={formik.values.captchaInput}
          onChange={(val) => formik.setFieldValue('captchaInput', val)}
        />
        {getFormErrorMessage('captchaInput')}
      </form>
    </Dialog>
  );
}
