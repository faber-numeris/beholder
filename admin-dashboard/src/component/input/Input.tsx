import {type Control, Controller, type FieldErrors, type FieldValues, type Path, type ControllerRenderProps } from 'react-hook-form';
import { InputText, type InputTextProps } from 'primereact/inputtext';
import { InputNumber, type InputNumberProps } from 'primereact/inputnumber';
import { Password, type PasswordProps } from 'primereact/password';
import { classNames } from 'primereact/utils';

type OverlapKeys = 'value' | 'onChange' | 'onBlur' | 'ref' | 'id' | 'name' | 'invalid';

type InputProps<T extends FieldValues = FieldValues> =
    & Omit<InputTextProps, OverlapKeys>
    & Omit<InputNumberProps, OverlapKeys>
    & Omit<PasswordProps, OverlapKeys>
    & {
        control: Control<T>;
        name: Path<T>;
        label: string;
        type?: 'text' | 'password' | 'number';
        rules?: Record<string, unknown>;
        errors?: FieldErrors<T>;
    };

export const Input = <T extends FieldValues>({
    control,
    name,
    label,
    type = 'text',
    rules,
    errors,
    className: userClassName,
    ...rest
}: InputProps<T>) => {
    const errorMessage = errors?.[name]?.message as string | undefined;
    const hasError = !!errorMessage;

    const renderInput = (
        field: ControllerRenderProps<T, Path<T>>,
        fieldState: { invalid: boolean }) => {

        const className = classNames(userClassName, { 'p-invalid': fieldState.invalid });
        switch (type) {
            case 'password':
                return <Password id={name} {...field} {...rest} toggleMask className={className} />;
            case 'number':
                return (
                    <InputNumber
                        id={name}
                        value={field.value as number}
                        onValueChange={(e) => field.onChange(e.value)}
                        {...rest}
                        className={className}
                    />
                );
            default:
                return <InputText id={name} {...field} {...rest} className={className} />;
        }
    };

    return (
        <div className="field">
            <span className="p-float-label">
                <Controller
                    name={name}
                    control={control}
                    rules={rules}
                    render={({ field, fieldState }) => renderInput(field, fieldState)}
                />
                <label htmlFor={name} className={classNames({ 'p-error': hasError })}>
                    {label}
                </label>
            </span>
            {hasError && <small className="p-error">{errorMessage}</small>}
        </div>
    );
};
