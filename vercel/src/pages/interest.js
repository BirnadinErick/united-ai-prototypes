import { supabase } from "../db";

export const prerender = false;

export const POST = async ({ request }) => {
  const data = await request.formData();

  const supabaseUrl = import.meta.env.SUPABASE_URL;
  const supabaseKey = import.meta.env.SUPABASE_KEY;

  const email = data.get("email");
  const course = data.get("course");
  const is_member = data.get("is_member") !== 1 ? 0 : 1;
  const project = data.get("project");

  const { error } = await supabase.from("interests").insert({
    email: email,
    course: course,
    is_member: is_member,
    project: project,
  });
  if (error) throw error;

  console.log(email);
  console.log(course);
  console.log(is_member);
  console.log(project);

  return new Response(`
<p class="my-6 bg-green-200 p-2 rounded-sm">
    Your interest was recorded. Further communication will be carried over through the
    email you provided.
</p>`);
};
