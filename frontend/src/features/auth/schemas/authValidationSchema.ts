import * as z from "zod";

export const loginSchema = z.object({
  email: z.string().email({ message: "Invalid email address." }),
  password: z.string().min(1, { message: "Password is required." }),
});

export type LoginFormValues = z.infer<typeof loginSchema>;

export const registrationSchema = z
  .object({
    email: z.string().email({ message: "Invalid email address." }),
    password: z
      .string()
      .min(8, { message: "Password must be at least 8 characters." }),
    confirmPassword: z.string(),
    name: z.string().min(1, { message: "Name is required." }),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords don't match.",
    path: ["confirmPassword"], // path of error
  });

export type RegistrationFormValues = z.infer<typeof registrationSchema>;
