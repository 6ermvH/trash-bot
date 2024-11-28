from telegram import Update, InlineKeyboardButton, InlineKeyboardMarkup
from telegram.ext import Application, CommandHandler, CallbackContext, CallbackQueryHandler

persons = ['Anthon', 'German', 'Vitalya', 'Seryoga', 'Effe']
current_person_index = 0
trash_taken_out_today = False

async def start(update: Update, context: CallbackContext) -> None:
    global trash_taken_out_today
    status_message = "Мусор уже вынесли сегодня." if trash_taken_out_today else "Мусор еще не вынесли сегодня."
    keyboard = [
        [InlineKeyboardButton("Кто выносит мусор", callback_data='1')],
        [InlineKeyboardButton("Мусор вынесен", callback_data='2')]
    ]
    reply_markup = InlineKeyboardMarkup(keyboard)
    await update.message.reply_text(f'{status_message}\nВыберите действие:', reply_markup=reply_markup)

async def button_handler(update: Update, context: CallbackContext) -> None:
    global current_person_index, trash_taken_out_today
    query = update.callback_query
    await query.answer()

    if query.data == '1':
        person = persons[current_person_index]
        await query.edit_message_text(text=f"Сегодня мусор выносит: {person}")
    elif query.data == '2':
        if trash_taken_out_today:
            await query.edit_message_text(text="Мусор уже вынесен сегодня. Чтобы вынести повторно, пропишите /next")
            return
        current_person_index = (current_person_index + 1) % len(persons)
        person = persons[current_person_index]
        trash_taken_out_today = True
        await query.edit_message_text(text=f"Теперь мусор выносит: {person}")

async def set_establish(update: Update, context: CallbackContext) -> None:
    global persons
    if context.args:
        # Очищаем старый список и добавляем новые имена
        persons = []
        persons.extend(context.args)
        await update.message.reply_text(f"Новый порядок установлен:\n{', '.join(persons)}")
    else:
        await update.message.reply_text("Пожалуйста, укажите список имен через пробел.\nПример: /set_establish Иван Петр Алексей")

async def next_day(update: Update, context: CallbackContext) -> None:
    global trash_taken_out_today
    trash_taken_out_today = False
    await update.message.reply_text("Флаг сброшен. Можно выносить мусор снова!")

async def prev_person(update: Update, context: CallbackContext) -> None:
    global current_person_index, trash_taken_out_today
    current_person_index = (current_person_index - 1) % len(persons)
    trash_taken_out_today = False
    person = persons[current_person_index]
    await update.message.reply_text(f"Возврат к предыдущему человеку: {person}\nФлаг сброшен.")

async def help_command(update: Update, context: CallbackContext) -> None:
    help_text = """
🗑 *Команды бота:*

/start \- Запустить бота и показать текущий статус
/help \- Показать это сообщение
/set\_establish \[имена\] \- Установить новый список людей \(через пробел\)
/next \- Сбросить флаг выноса мусора
/prev \- Вернуться к предыдущему человеку

*Кнопки:*
• Кто выносит мусор \- Показать, кто сейчас должен выносить мусор
• Мусор вынесен \- Отметить, что мусор вынесен, и перейти к следующему человеку
"""
    await update.message.reply_text(help_text, parse_mode='MarkdownV2')

def main() -> None:
    application = Application.builder().token("7853464150:AAG-hYnlKSHv9zrMPIDqtlfv0MoL1rQ_PI4").build()

    application.add_handler(CommandHandler("start", start))
    application.add_handler(CommandHandler("help", help_command))
    application.add_handler(CommandHandler("set_establish", set_establish))
    application.add_handler(CommandHandler("next", next_day))
    application.add_handler(CommandHandler("prev", prev_person))
    application.add_handler(CallbackQueryHandler(button_handler))

    print("Бот запущен...")
    application.run_polling(allowed_updates=Update.ALL_TYPES, timeout=60)
    print("Бот остановлен...")

if __name__ == '__main__':
    import asyncio
    asyncio.run(main())
