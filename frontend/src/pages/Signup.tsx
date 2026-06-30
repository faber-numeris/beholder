import React, {useState} from 'react';
import {useForm} from 'react-hook-form';
import {Button} from 'primereact/button';
import {Fieldset} from 'primereact/fieldset'
import {Input} from '../component/input/Input';

interface SignupForm {
    email: string;
    password: string;
    passwordConfirmation: string
}

const defaultValues: SignupForm = {
    email: '',
    password: '',
    passwordConfirmation: ''
}

export const Signup: React.FC = () => {
    const [formData, setFormData] = useState<SignupForm>();

    const {control, formState: {errors}, handleSubmit, reset} = useForm({defaultValues});

    const onSubmit = (data: SignupForm) => {
        setFormData(data);
        console.log(formData)
        reset();
    };

    return (
        <div className="flex align-items-center justify-content-center min-h-screen p-4">
            <div className="w-full lg:w-7">
                <Fieldset legend={"Register"}>
                    <form onSubmit={handleSubmit(onSubmit)} className="p-fluid gap-1 flex flex-column">
                        <Input control={control}
                               name="email"
                               label="email*"
                               type="text"
                               rules={{required: 'Email is required.'}}
                               errors={errors}/>
                        <Input control={control}
                               name="password"
                               label="password*"
                               type="password"
                               rules={{required: 'Password is required.'}}
                               errors={errors}
                        />
                        <Input control={control}
                               name="passwordConfirmation"
                               label="password check*"
                               type="password"
                               rules={{required: 'Password confirmation is required.'}}
                               errors={errors}/>
                        <Button
                            type="submit"
                            label="Submit"
                            className="mt-2"/>
                    </form>
                </Fieldset>
            </div>
        </div>
    );
}

export default Signup
